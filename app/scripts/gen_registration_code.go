package main

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// gen_registration_code creates a bootstrap invite code when subscriptions are disabled.
func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	if !config.RegistrationGateEnabled() {
		fmt.Fprintln(os.Stderr, "registration codes require SUBSCRIPTIONS_ENABLED=false")
		os.Exit(1)
	}

	dbPath := config.EnvOr("DATABASE_PATH", "./database/main.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	store := storage.NewSqliteStorage(db)
	authService := auth.NewAuthService(store, nil)

	code, expiresAt, err := authService.GenerateRegistrationCode(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate registration code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Registration code: %s\n", code)
	fmt.Printf("Expires at: %s\n", time.Unix(expiresAt, 0).UTC().Format(time.RFC3339))
}
