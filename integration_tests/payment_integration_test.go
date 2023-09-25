package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"payment-api/internal/config"
	"payment-api/internal/db"
	"payment-api/internal/logger"
	"payment-api/internal/models"
	"payment-api/internal/server"
)

type PaymentTestSuite struct {
	suite.Suite
	dbConn *sql.DB
	cnf    *config.Config
}

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, &PaymentTestSuite{})
}

func (s *PaymentTestSuite) SetupSuite() {
	os.Setenv("DB_HOST", "0.0.0.0")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("PROVIDER_FILE_PATH", "../assets/providers.json")
	os.Setenv("STORES_FILE_PATH", "../assets/stores.json")

	cnf := config.Load()
	s.cnf = cnf
	lg, err := logger.New(cnf.Environment, cnf.LogLevel)
	if err != nil {
		log.Fatalf("failed to initialize a logger, err: %v", err)
	}
	dbConn, err := db.Open(context.Background(), cnf)
	if err != nil {
		lg.Fatalf("failed to open db connection")
	}

	s.dbConn = dbConn
	err = db.InitialSeed(dbConn, []string{db.CreateProvider, db.SeedProviderApplePay, db.SeedProviderGooglePay, db.SeedProviderPayPal, db.SeedProviderStripe, db.SeedInvalidProvider}, lg)
	if err != nil {
		lg.Fatalf("failed to execute statements, error: %v", err)
	}
	// Launching the server
	go func() {
		server.Run(lg, cnf, dbConn)
	}()
}

// TestApp checks if the endpoint works well, from request to response
func (s *PaymentTestSuite) TestApp() {
	res, err := s.dbConn.Query("SELECT id, name FROM providers ORDER BY created_at")
	if err != nil {
		s.FailNow("failed to query rows from providers")
	}
	var providers []models.Provider
	for res.Next() {
		var p models.Provider
		if err := res.Scan(&p.ID, &p.Name); err != nil {
			s.FailNow("failed to scan columns from providers")
		}
		providers = append(providers, p)
	}
	type testCase struct {
		name string
		id   string
		data string
		urls []map[string]string
		code int
		msg  string
	}

	// TODO: add abscent id case
	testCases := []testCase{
		{
			name: "success ApplePay",
			id:   providers[0].ID,
			data: "https://apple-pay-gateway.apple.com",
			urls: nil,
			code: http.StatusOK,
			msg:  "",
		},
		{
			name: "success GooglePay",
			id:   providers[1].ID,
			data: "https://google.com/pay",
			urls: nil,
			code: http.StatusOK,
			msg:  "",
		},
		{
			name: "success PayPal",
			id:   providers[2].ID,
			data: "https://www.paypal.com/pay",
			urls: nil,
			code: http.StatusOK,
			msg:  "",
		},
		{
			name: "success Stripe",
			id:   providers[3].ID,
			data: "https://stripe.com/pay",
			urls: nil,
			code: http.StatusOK,
			msg:  "",
		},
		{
			name: "success Stores",
			id:   providers[4].ID,
			data: "",
			urls: []map[string]string{
				{
					"AppleStore": "https://apps.apple.com/us/app/headway-daily-book-summaries/id1457185832",
				},
				{
					"PlayMarket": "https://play.google.com/store/apps/details?id=com.headway.books&hl=pl&pli=1",
				},
			},
			code: http.StatusOK,
			msg:  "",
		},
		{
			name: "fail nonexistent provider",
			id:   "7626be3d-06ea-43d0-895c-dfbf017c7fff",
			data: "",
			urls: nil,
			code: http.StatusBadRequest,
			msg:  "Provided parameter has bad format",
		},
		{
			name: "fail no parameter",
			id:   "",
			data: "",
			urls: nil,
			code: http.StatusBadRequest,
			msg:  "Missing productID parameter",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := http.Get(
				"http://0.0.0.0:" + s.cnf.Service.Port +
					"/api/v1/payment/url?productID=" + tc.id)
			if err != nil {
				s.FailNow("failed to perform a request to the server")
			}
			defer resp.Body.Close()
			type respMsg struct {
				Code    int
				Data    string
				Urls    []map[string]string `json:"stores_urls"`
				Message string              `json:"message"`
			}

			var respData respMsg
			_ = json.NewDecoder(resp.Body).Decode(&respData)
			s.Equal(tc.code, respData.Code)
			s.Equal(tc.data, respData.Data)
			s.Equal(tc.msg, respData.Message)

			// If stores are returned
			if tc.urls != nil {
				s.Equal(tc.urls[0]["AppleStore"], respData.Urls[0]["AppleStore"])
				s.Equal(tc.urls[1]["PlayMarket"], respData.Urls[1]["PlayMarket"])
			}
		})
	}
}
