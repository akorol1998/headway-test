package main

import (
	"context"
	"log"

	"payment-api/internal/config"
	"payment-api/internal/db"
	"payment-api/internal/logger"
	"payment-api/internal/server"
)

func main() {
	cnf := config.Load()
	lg, err := logger.New(cnf.Environment, cnf.LogLevel)
	if err != nil {
		log.Fatalf("failed to initialize a logger, err: %v", err)
	}
	dbConn, err := db.Open(context.Background(), cnf)
	if err != nil {
		lg.Fatalf("failed to open db connection")
	}
	err = db.InitialSeed(dbConn, []string{db.CreateProvider, db.SeedProviderApplePay, db.SeedProviderGooglePay, db.SeedProviderPayPal, db.SeedProviderStripe, db.SeedInvalidProvider}, lg)
	if err != nil {
		lg.Fatalf("failed to execute statements, error: %v", err)
	}
	// Log providers
	db.Providers(dbConn, lg)

	server.Run(lg, cnf, dbConn)
}
