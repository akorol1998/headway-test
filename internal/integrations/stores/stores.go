package stores

import (
	"encoding/json"
	"errors"
	"os"

	"go.uber.org/zap"
)

const (
	StoreAppleStore = "AppleStore"
	StorePlayMarket = "PlayMarket"
)

var ErrUnknownStoreName = errors.New("unknown store name")

type StoreName string

type configStores struct {
	PlayMarket string `json:"play_market"`
	AppleStore string `json:"apple_store"`
}

type Stores struct {
	log      *zap.SugaredLogger
	filePath string
}

func NewStore(log *zap.SugaredLogger, filePath string) *Stores {
	return &Stores{log: log, filePath: filePath}
}

// AppUrl fetches url to the Headway app from the provided store name
// This method is coupled to business logic, thus is tested within business logic scope
func (s *Stores) AppUrl(name StoreName) (string, error) {
	raw, err := os.ReadFile(s.filePath)
	if err != nil {
		s.log.Errorf("failed to read the file, error: %v", err)
		return "", err
	}

	var cnf configStores
	if err := json.Unmarshal(raw, &cnf); err != nil {
		s.log.Errorf("failed to unmarshall stores config file, error: %v", err)
		return "", err
	}
	switch name {
	case StoreAppleStore:
		return cnf.AppleStore, nil
	case StorePlayMarket:
		return cnf.PlayMarket, nil
	default:
		return "", ErrUnknownStoreName
	}
}
