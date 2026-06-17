package main

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/server"
	"FinancialTracker/internal/storage"
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Unable to load .env file")
		}
	}
    // Opens the test_schema.sql file
    schemaFile, err := os.ReadFile("./database/test_schema.sql")
    if err != nil {
        log.Fatal("Missing test_schema.sql file!")
    }

    db, err := sql.Open("sqlite3", "./database/main.db?_foreign_keys=on")
    if err != nil {
        log.Fatal("Unable to open main.db!")
    }
    defer db.Close()

    // Execute schema
    if _, err := db.Exec(string(schemaFile)); err != nil {
        log.Fatal("Unable to execute schema sql:", err)
    }
	log.Debug("TETING")
	store := storage.NewSqliteStorage(db)

	// **************** Start of creating test account ************************
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	
	user := &models.User{
	    Email:           "test@test.com",
	    FirstName:       "test",
	    LastName:        "test",
	    PasswordHash:    string(hash),
	    ThemePreference: "system",
	}
	
	if err := store.CreateUser(user); err != nil {
		log.Error(err)
	}
	// *************** End of creating test account *********************

	// Run the api server
	e := server.RunServer(store)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	<-ctx.Done()
}
