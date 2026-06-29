package auth

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/mail"
	"FinancialTracker/internal/storage"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type recordingMailSender struct {
	lastTo        string
	lastFirstName string
	lastCode      string
}

func (r *recordingMailSender) SendPasswordResetCode(to, firstName, code string) error {
	r.lastTo = to
	r.lastFirstName = firstName
	r.lastCode = code
	return nil
}

func openAuthTestDB(t *testing.T) *storage.Storage {
	t.Helper()

	_, file, _, _ := runtime.Caller(0)
	schemaPath := filepath.Join(filepath.Dir(file), "../../../database/test_schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read test schema: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	return storage.NewSqliteStorage(db)
}

func createPasswordUser(t *testing.T, store *storage.Storage, email, password string) *models.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &models.User{
		Email:        email,
		FirstName:    "Test",
		LastName:     "User",
		PasswordHash: string(hash),
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func TestRequestPasswordReset_unknownEmailIsEnumerationSafe(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	if err := svc.RequestPasswordReset("missing@example.com"); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if sender.lastCode != "" {
		t.Fatal("expected no mail for unknown email")
	}
}

func TestRequestPasswordReset_resendCooldownAllowsFirstResend(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "cooldown@example.com", "OldPass!1")
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("initial RequestPasswordReset() error = %v", err)
	}
	firstCode := sender.lastCode
	if firstCode == "" {
		t.Fatal("expected initial reset code")
	}

	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("first resend RequestPasswordReset() error = %v", err)
	}
	if sender.lastCode == firstCode {
		t.Fatal("expected first resend to issue a new code")
	}

	secondResendCode := sender.lastCode
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("second resend RequestPasswordReset() error = %v", err)
	}
	if sender.lastCode != secondResendCode {
		t.Fatal("expected second resend within cooldown to be ignored")
	}
}

func TestRequestPasswordReset_skipsSSOOnlyAccount(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := &models.User{
		Email:     "sso@example.com",
		FirstName: "SSO",
		LastName:  "User",
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if sender.lastCode != "" {
		t.Fatal("expected no mail for SSO-only account")
	}
}

func TestConfirmPasswordReset_success(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "reset@example.com", "OldPass!1")

	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	if sender.lastCode == "" {
		t.Fatal("expected reset code to be issued")
	}

	if err := svc.ConfirmPasswordReset(user.Email, sender.lastCode, "NewPass!2", "NewPass!2"); err != nil {
		t.Fatalf("ConfirmPasswordReset() error = %v", err)
	}

	updated, err := store.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail() error = %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("NewPass!2")); err != nil {
		t.Fatal("password was not updated")
	}

	if _, err := svc.Authenticate(user.Email, "NewPass!2", false); err != nil {
		t.Fatalf("Authenticate() with new password error = %v", err)
	}
}

func TestVerifyPasswordResetCode_success(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "verify@example.com", "OldPass!1")
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	expiresAt, err := svc.VerifyPasswordResetCode(user.Email, sender.lastCode)
	if err != nil {
		t.Fatalf("VerifyPasswordResetCode() error = %v", err)
	}
	if expiresAt <= time.Now().Unix() {
		t.Fatalf("expiresAt = %d, want future timestamp", expiresAt)
	}

	if err := svc.ConfirmPasswordReset(user.Email, sender.lastCode, "NewPass!2", "NewPass!2"); err != nil {
		t.Fatalf("ConfirmPasswordReset() after verify error = %v", err)
	}
}

func TestConfirmPasswordReset_wrongCodeIncrementsAttempts(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "wrong@example.com", "OldPass!1")
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	for i := 0; i < resetMaxAttempts; i++ {
		err := svc.ConfirmPasswordReset(user.Email, "000000", "NewPass!2", "NewPass!2")
		if err != ErrInvalidResetCode {
			t.Fatalf("attempt %d: error = %v, want ErrInvalidResetCode", i+1, err)
		}
	}

	err := svc.ConfirmPasswordReset(user.Email, sender.lastCode, "NewPass!2", "NewPass!2")
	if err != ErrInvalidResetCode {
		t.Fatalf("after max attempts: error = %v, want ErrInvalidResetCode", err)
	}
}

func TestVerifyPasswordResetCode_wrongCodeIncrementsAttempts(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "verify-wrong@example.com", "OldPass!1")
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	_, err := svc.VerifyPasswordResetCode(user.Email, "000000")
	if err != ErrInvalidResetCode {
		t.Fatalf("VerifyPasswordResetCode() error = %v, want ErrInvalidResetCode", err)
	}

	if err := svc.ConfirmPasswordReset(user.Email, sender.lastCode, "NewPass!2", "NewPass!2"); err != nil {
		t.Fatalf("ConfirmPasswordReset() with correct code after one failed verify error = %v", err)
	}
}

func TestConfirmPasswordReset_expiredCodeRejected(t *testing.T) {
	store := openAuthTestDB(t)
	sender := &recordingMailSender{}
	svc := NewAuthService(store, sender)

	user := createPasswordUser(t, store, "expired@example.com", "OldPass!1")
	if err := svc.RequestPasswordReset(user.Email); err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(sender.lastCode), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash code: %v", err)
	}
	if err := store.InvalidatePasswordResetCodes(user.ID); err != nil {
		t.Fatalf("InvalidatePasswordResetCodes() error = %v", err)
	}
	expiredAt := time.Now().Add(-time.Minute).Unix()
	if err := store.CreatePasswordResetCode(user.ID, string(hash), expiredAt); err != nil {
		t.Fatalf("CreatePasswordResetCode() error = %v", err)
	}

	err = svc.ConfirmPasswordReset(user.Email, sender.lastCode, "NewPass!2", "NewPass!2")
	if err != ErrInvalidResetCode {
		t.Fatalf("ConfirmPasswordReset() error = %v, want ErrInvalidResetCode", err)
	}
}

func TestConfirmPasswordReset_invalidCodeFormat(t *testing.T) {
	store := openAuthTestDB(t)
	svc := NewAuthService(store, mail.DevSender{})

	err := svc.ConfirmPasswordReset("reset@example.com", "abc", "NewPass!2", "NewPass!2")
	if err == nil || err.Error() != "reset code must be 6 digits" {
		t.Fatalf("ConfirmPasswordReset() error = %v", err)
	}
}

func TestGenerateResetCode_format(t *testing.T) {
	code, err := generateResetCode()
	if err != nil {
		t.Fatalf("generateResetCode() error = %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("code length = %d, want 6", len(code))
	}
	for _, ch := range code {
		if ch < '0' || ch > '9' {
			t.Fatalf("code contains non-digit: %q", code)
		}
	}
}
