package main

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/server"
	"FinancialTracker/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	loadEnvFiles()

	if err := config.ValidateRequiredEnv(); err != nil {
		log.Fatal(err)
	}

	db, store, err := openDatabase(config.EnvOr("DATABASE_PATH", "./database/main.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if config.IsDevelopment() {
		seedDevelopmentUser(store)
	}

	e := server.RunServer(store)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":" + config.EnvOr("PORT", "8080")
	sc := echo.StartConfig{
		Address:         addr,
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	<-ctx.Done()
}

// loadEnvFiles loads a local .env file when present. Production should inject env vars directly.
func loadEnvFiles() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
}

// openDatabase opens SQLite and applies the schema appropriate for the current environment.
func openDatabase(path string) (*sql.DB, *storage.Storage, error) {
	schemaPath := "./database/schema.sql"
	if config.IsDevelopment() {
		schemaPath = "./database/test_schema.sql"
	}

	schemaFile, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read schema %s: %w", schemaPath, err)
	}

	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if _, err := db.Exec(string(schemaFile)); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("apply schema: %w", err)
	}

	return db, storage.NewSqliteStorage(db), nil
}

// seedDevelopmentUser creates a default login for local testing when no account exists yet.
func seedDevelopmentUser(store *storage.Storage) {
	existing, err := store.GetUserByEmail("test@test.com")
	if err != nil {
		log.Errorf("development seed lookup failed: %v", err)
		return
	}
	if existing != nil {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("development seed password hash failed: %v", err)
		return
	}

	user := &models.User{
		Email:           "test@test.com",
		FirstName:       "test",
		LastName:        "test",
		PasswordHash:    string(hash),
		ThemePreference: "system",
	}
	if err := store.CreateUser(user); err != nil {
		log.Errorf("development seed user failed: %v", err)
	}
}
