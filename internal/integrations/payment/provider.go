package payment

import (
	"encoding/json"
	"errors"
	"os"

	"go.uber.org/zap"

	"payment-api/internal/models"
)

var (
	ErrUnknownProviderID = errors.New("unknown provider name")
)

type configProviders struct {
	ApplePay  string `json:"apple_pay"`
	GooglePay string `json:"google_pay"`
	StripPay  string `json:"stripe"`
	PayPal    string `json:"pay_pal"`
}

type PaymentProvider struct {
	log      *zap.SugaredLogger
	filePath string
}

func NewPaymentProvider(log *zap.SugaredLogger, filePath string) *PaymentProvider {
	return &PaymentProvider{log: log, filePath: filePath}
}

// PaymentUrl mocks process of generating link for the provider which name was passed method
// since it is a mock which is coupled to business logic, thus is tested within it
func (p *PaymentProvider) PaymentUrl(name, apiKey, secret string) (string, error) {
	raw, err := os.ReadFile(p.filePath)
	if err != nil {
		p.log.Errorf("failed to read the file, error: %v", err)
		return "", err
	}

	var cnf configProviders
	if err := json.Unmarshal(raw, &cnf); err != nil {
		p.log.Errorf("failed to unmarshall providers config file, error: %v", err)
		return "", err
	}

	p.log.Infof("paymentProvider: generating a link for: %v", name)
	switch name {
	case models.ProviderNameApplePay:
		return cnf.ApplePay, nil
	case models.ProviderNameGooglePay:
		return cnf.GooglePay, nil
	case models.ProviderNamePayPal:
		return cnf.PayPal, nil
	case models.ProviderNameStripe:
		return cnf.StripPay, nil
	default:
		return "", ErrUnknownProviderID
	}
}
