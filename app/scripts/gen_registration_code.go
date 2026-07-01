package main

import (
	"FinancialTracker/internal/cmd"

	"github.com/joho/godotenv"
)

// gen_registration_code creates a bootstrap invite code when subscriptions are disabled.
func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	cmd.GenRegistrationCodeOrExit()
}
