package payment

import (
	"context"
	"fmt"
	"payment-api/internal/integrations/payment"
	"payment-api/internal/integrations/stores"
	"payment-api/internal/models"
	"payment-api/internal/services/payment/repository"
	"testing"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// FakeProviderRepo is faked structure for existing repository
type FakeProviderRepo struct {
	Providers []*models.Provider
}

func (m *FakeProviderRepo) Setup() {
	seeds := []map[string]string{
		{
			"id":   "f3694470-79e5-46c2-bcab-cad4750cdbcc",
			"name": "ApplePay",
		},
		{
			"id":   "3754705b-842a-4d53-a7e5-cbf8025a5ddb",
			"name": "GooglePay",
		},
		{
			"id":   "b6b37fd9-5cd9-4e07-8975-08e42a09f723",
			"name": "PayPal",
		},
		{
			"id":   "39251d76-1b3c-470d-969d-c7dade716d97",
			"name": "Stripe",
		},
		{
			"id":   "22881d76-1b3c-470d-969d-c7dade716d97",
			"name": "InvalidProvider",
		},
	}
	m.Providers = make([]*models.Provider, 0, len(seeds))
	for i, seed := range seeds {
		m.Providers = append(m.Providers,
			&models.Provider{
				ID:     seed["id"],
				Name:   seed["name"],
				ApiKey: fmt.Sprintf("api_key_%v", i),
				Secret: fmt.Sprintf("secret_%v", i),
			})
	}
}

func (m *FakeProviderRepo) FetchByID(id string) (*models.Provider, error) {
	for _, p := range m.Providers {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, repository.ErrNotFound
}

func TestPaymentServicePaymentUrl(t *testing.T) {
	mockLogger := zap.NewNop().Sugar()
	// Fake repo
	fakeProviderRepo := FakeProviderRepo{}
	fakeProviderRepo.Setup()
	// Since it already acts as a fake structure for mocking requests to the payment platforms
	// it will be used as it is
	paymentProvider := payment.NewPaymentProvider(mockLogger, "../../../assets/providers.json")

	service := NewPaymentService(mockLogger, paymentProvider, nil, &fakeProviderRepo)

	type testCase struct {
		name        string
		id          string
		expectedUrl string
		success     bool
		expectedErr error
	}
	// Defining test data
	testCases := []testCase{
		{
			"success Apple",
			fakeProviderRepo.Providers[0].ID,
			"https://apple-pay-gateway.apple.com",
			true,
			nil,
		},
		{
			"success Google",
			fakeProviderRepo.Providers[1].ID,
			"https://google.com/pay",
			true,
			nil,
		},
		{
			"success PayPal",
			fakeProviderRepo.Providers[2].ID,
			"https://www.paypal.com/pay",
			true,
			nil,
		},
		{
			"success Stripe",
			fakeProviderRepo.Providers[3].ID,
			"https://stripe.com/pay",
			true,
			nil,
		},
		{
			"fail provider",
			fakeProviderRepo.Providers[4].ID,
			"",
			false,
			ErrProvider,
		},
		{
			"fail non-existent ID",
			uuid.NewString(),
			"",
			false,
			ErrNotFound,
		},
		{
			"fail invalid id format",
			"ab23bd-123efa4b1",
			"",
			false,
			ErrUuidInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := service.PaymentUrl(context.Background(), tc.id)
			if tc.success {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
			assert.Equal(t, tc.expectedUrl, url)
		})
	}
}

func TestPaymentServiceStoresUrls(t *testing.T) {
	mockLogger := zap.NewNop().Sugar()
	// Since it already acts as a fake structure for mocking requests to the payment platforms
	// it will be used as it is
	paymentProvider := payment.NewPaymentProvider(mockLogger, "../../../assets/providers.json")
	// Initiating store dependency
	stores := stores.NewStore(mockLogger, "../../../assets/stores.json")
	service := NewPaymentService(mockLogger, paymentProvider, stores, nil)
	urlMap, err := service.StoresUrls(context.Background())

	assert.NoError(t, err)
	url, ok := urlMap[0]["AppleStore"]
	assert.True(t, ok)
	assert.Equal(t, url, "https://apps.apple.com/us/app/headway-daily-book-summaries/id1457185832")

	url, ok = urlMap[1]["PlayMarket"]
	assert.True(t, ok)
	assert.Equal(t, url, "https://play.google.com/store/apps/details?id=com.headway.books&hl=pl&pli=1")
}
