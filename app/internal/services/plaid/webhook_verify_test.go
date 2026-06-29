package plaid

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/plaid/plaid-go/v42/plaid"
)

func TestSha256Hex(t *testing.T) {
	body := []byte(`{"webhook_type":"TRANSACTIONS","webhook_code":"SYNC_UPDATES_AVAILABLE","item_id":"abc"}`)
	got := sha256Hex(body)
	if len(got) != 64 {
		t.Fatalf("sha256Hex() length = %d, want 64", len(got))
	}

	sum := sha256.Sum256(body)
	want := fmt.Sprintf("%x", sum[:])
	if got != want {
		t.Fatalf("sha256Hex() = %q, want %q", got, want)
	}
}

func TestDecodeBase64URL(t *testing.T) {
	decoded, err := decodeBase64URL("AQ")
	if err != nil {
		t.Fatalf("decodeBase64URL() error = %v", err)
	}
	if len(decoded) != 1 || decoded[0] != 0x01 {
		t.Fatalf("decodeBase64URL() = %v, want [1]", decoded)
	}
}

func TestJwkToECDSAPublicKeyRejectsUnsupportedCurve(t *testing.T) {
	jwk := plaid.JWKPublicKey{
		Kty: "RSA",
		Crv: "P-256",
		X:   "AQ",
		Y:   "Ag",
	}
	if _, err := jwkToECDSAPublicKey(jwk); err == nil {
		t.Fatal("expected error for non-EC key")
	}
}
