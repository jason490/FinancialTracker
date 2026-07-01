package cmd

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// GenRegistrationCode creates a bootstrap invite code when subscriptions are disabled.
func GenRegistrationCode() error {
	if !config.RegistrationGateEnabled() {
		return fmt.Errorf("registration codes require SUBSCRIPTIONS_ENABLED=false")
	}

	dbPath := config.EnvOr("DATABASE_PATH", "./database/main.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	store := storage.NewSqliteStorage(db)
	authService := auth.NewAuthService(store, nil)

	code, expiresAt, err := authService.GenerateRegistrationCode(nil)
	if err != nil {
		return fmt.Errorf("generate registration code: %w", err)
	}

	fmt.Printf("Registration code: %s\n", code)
	fmt.Printf("Expires at: %s\n", time.Unix(expiresAt, 0).UTC().Format(time.RFC3339))
	return nil
}

// GenRegistrationCodeOrExit runs GenRegistrationCode and exits the process on failure.
func GenRegistrationCodeOrExit() {
	if err := GenRegistrationCode(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
