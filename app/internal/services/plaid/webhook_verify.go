package plaid

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plaid/plaid-go/v42/plaid"
)

const webhookMaxAge = 5 * time.Minute

var (
	ErrWebhookVerification = errors.New("webhook verification failed")
)

type webhookKeyCache struct {
	mu  sync.RWMutex
	kid string
	key plaid.JWKPublicKey
}

// verifyWebhook validates the Plaid-Verification JWT and request body hash.
func (p *PlaidService) verifyWebhook(ctx context.Context, body []byte, signedJWT string) error {
	signedJWT = strings.TrimSpace(signedJWT)
	if signedJWT == "" {
		return ErrWebhookVerification
	}

	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(signedJWT, jwt.MapClaims{})
	if err != nil {
		return ErrWebhookVerification
	}

	alg, _ := token.Header["alg"].(string)
	if alg != jwt.SigningMethodES256.Alg() {
		return ErrWebhookVerification
	}

	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return ErrWebhookVerification
	}

	jwk, err := p.getWebhookVerificationKey(ctx, kid)
	if err != nil {
		return ErrWebhookVerification
	}

	publicKey, err := jwkToECDSAPublicKey(jwk)
	if err != nil {
		return ErrWebhookVerification
	}

	parsed, err := jwt.Parse(signedJWT, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return publicKey, nil
	})
	if err != nil || !parsed.Valid {
		return ErrWebhookVerification
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return ErrWebhookVerification
	}

	iat, ok := claims["iat"].(float64)
	if !ok || time.Now().Unix()-int64(iat) > int64(webhookMaxAge.Seconds()) {
		return ErrWebhookVerification
	}

	expectedHash, ok := claims["request_body_sha256"].(string)
	if !ok || expectedHash == "" {
		return ErrWebhookVerification
	}

	actualHash := sha256Hex(body)
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(actualHash)) != 0 {
		return ErrWebhookVerification
	}

	return nil
}

// getWebhookVerificationKey loads and caches the JWK used to verify Plaid webhooks.
func (p *PlaidService) getWebhookVerificationKey(ctx context.Context, keyID string) (plaid.JWKPublicKey, error) {
	p.webhookKeyCache.mu.RLock()
	if p.webhookKeyCache.kid == keyID {
		key := p.webhookKeyCache.key
		p.webhookKeyCache.mu.RUnlock()
		return key, nil
	}
	p.webhookKeyCache.mu.Unlock()

	request := plaid.NewWebhookVerificationKeyGetRequest(keyID)
	resp, _, err := p.client.PlaidApi.WebhookVerificationKeyGet(ctx).WebhookVerificationKeyGetRequest(*request).Execute()
	if err != nil {
		return plaid.JWKPublicKey{}, err
	}

	key := resp.GetKey()
	p.webhookKeyCache.mu.Lock()
	p.webhookKeyCache.kid = keyID
	p.webhookKeyCache.key = key
	p.webhookKeyCache.mu.Unlock()

	return key, nil
}

// jwkToECDSAPublicKey converts a Plaid JWK into an ECDSA public key for ES256 verification.
func jwkToECDSAPublicKey(jwk plaid.JWKPublicKey) (*ecdsa.PublicKey, error) {
	if jwk.GetKty() != "EC" || jwk.GetCrv() != "P-256" {
		return nil, fmt.Errorf("unsupported jwk")
	}

	xBytes, err := decodeBase64URL(jwk.GetX())
	if err != nil {
		return nil, err
	}
	yBytes, err := decodeBase64URL(jwk.GetY())
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:    new(big.Int).SetBytes(xBytes),
		Y:    new(big.Int).SetBytes(yBytes),
	}, nil
}

// decodeBase64URL decodes a base64url-encoded string with optional padding.
func decodeBase64URL(value string) ([]byte, error) {
	if value == "" {
		return nil, fmt.Errorf("empty base64url value")
	}
	switch len(value) % 4 {
	case 2:
		value += "=="
	case 3:
		value += "="
	}
	return base64.URLEncoding.DecodeString(value)
}

// sha256Hex returns the lowercase hex-encoded SHA-256 digest of body.
func sha256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return fmt.Sprintf("%x", sum[:])
}

// webhookURL resolves the callback URL Plaid should send webhooks to.
func (p *PlaidService) webhookURL() string {
	if url := strings.TrimSpace(os.Getenv("PLAID_WEBHOOK_URL")); url != "" {
		return strings.TrimRight(url, "/")
	}
	if base := strings.TrimSpace(os.Getenv("API_PUBLIC_URL")); base != "" {
		return strings.TrimRight(base, "/") + "/api/v1/plaid/webhook"
	}
	return ""
}
