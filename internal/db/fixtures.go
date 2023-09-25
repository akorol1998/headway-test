package db

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	CreateProvider = `
	CREATE TABLE providers(
		id UUID PRIMARY KEY,
		name VARCHAR(32),
		api_key VARCHAR(255),
		secret VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	SeedProviderPayPal = fmt.Sprintf(`
	INSERT INTO providers (id,name,api_key,secret) VALUES ('%v', '%v', '%v', '%v');
	`, uuid.New().String(), "PayPal", "test_api_key", "test_secret")
	SeedProviderApplePay = fmt.Sprintf(`
	INSERT INTO providers (id,name,api_key,secret) VALUES ('%v', '%v', '%v', '%v');
	`, uuid.New().String(), "ApplePay", "test_api_key", "test_secret")
	SeedProviderGooglePay = fmt.Sprintf(`
	INSERT INTO providers (id,name,api_key,secret) VALUES ('%v', '%v', '%v', '%v');
	`, uuid.New().String(), "GooglePay", "test_api_key", "test_secret")
	SeedProviderStripe = fmt.Sprintf(`
	INSERT INTO providers (id,name,api_key,secret) VALUES ('%v', '%v', '%v', '%v');
	`, uuid.New().String(), "Stripe", "test_api_key", "test_secret")
	SeedInvalidProvider = fmt.Sprintf(`
	INSERT INTO providers (id,name,api_key,secret) VALUES ('%v', '%v', '%v', '%v');
	`, uuid.New().String(), "InvalidProvider", "test_api_key", "test_secret")
)
