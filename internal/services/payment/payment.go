package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"payment-api/internal/integrations/stores"
	"payment-api/internal/models"
	"payment-api/internal/services/payment/repository"
)

type Stores interface {
	// AppUrl retrieves url to the app from a provided store
	AppUrl(name stores.StoreName) (string, error)
}

// PaymentProvider allows you to work with provider implementation
type PaymentProvider interface {
	// PaymentUrl fetches url that is needed for payment
	PaymentUrl(name, apiKey, secret string) (string, error)
}

// Repository for provider
type ProviderRepo interface {
	FetchByID(id string) (*models.Provider, error)
}

type PaymentService struct {
	log             *zap.SugaredLogger
	paymentProvider PaymentProvider
	stores          Stores
	providerRepo    ProviderRepo
}

func NewPaymentService(log *zap.SugaredLogger, paymentProvider PaymentProvider, stores Stores, providerRepo ProviderRepo) *PaymentService {
	return &PaymentService{log: log, paymentProvider: paymentProvider, stores: stores, providerRepo: providerRepo}
}

// PaymentUrl returns payment url for the provided providerID
func (s *PaymentService) PaymentUrl(ctx context.Context, providerID string) (string, error) {
	// validating a uuid, since this logic may be used from more than one handler
	_, err := uuid.Parse(providerID)
	if err != nil {
		s.log.Errorw("failed to validate providerID",
			"ID", providerID)
		return "", ErrUuidInvalidFormat
	}

	// not parsing context for the sake of simplicity of the case
	providerModel, err := s.providerRepo.FetchByID(providerID)
	if err != nil {
		s.log.Errorw("failed to fetch provider by ID",
			"ID", providerID)
		switch err {
		case repository.ErrNotFound:
			return "", ErrNotFound
		case repository.ErrUuidInvalidFormat:
			return "", ErrUuidInvalidFormat
		default:
			return "", ErrUuidInvalidFormat
		}
	}

	// Instead of name could be used ENUM enumeration in the form of iota
	url, err := s.paymentProvider.PaymentUrl(providerModel.Name, providerModel.ApiKey, providerModel.Secret)
	if err != nil {
		s.log.Errorf("failed to get url from %v provider, error: %v", providerModel.Name, err)
		return "", ErrProvider
	}
	return url, nil
}

// StoresUrls fetches urls to all available stores where app is hosted
func (s *PaymentService) StoresUrls(ctx context.Context) ([]map[string]string, error) {
	urls := make([]map[string]string, 0, 2)
	for _, i := range []stores.StoreName{stores.StoreAppleStore, stores.StorePlayMarket} {
		url, err := s.stores.AppUrl(i)
		if err != nil {
			return nil, ErrStore
		}
		// not the best solution, though for such case it is much more convennient
		// than using fixed size arrays
		urls = append(urls, map[string]string{string(i): url})
	}
	return urls, nil
}
