package models

import "time"

const (
	ProviderNameApplePay  = "ApplePay"
	ProviderNameGooglePay = "GooglePay"
	ProviderNamePayPal    = "PayPal"
	ProviderNameStripe    = "Stripe"
)

type Provider struct {
	ID        string
	Name      string
	ApiKey    string
	Secret    string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
