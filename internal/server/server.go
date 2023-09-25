package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"payment-api/internal/config"
	intpayment "payment-api/internal/integrations/payment"
	"payment-api/internal/integrations/stores"
	"payment-api/internal/middlwares"
	"payment-api/internal/services/payment"
	v1 "payment-api/internal/services/payment/handlers/http/v1"
	"payment-api/internal/services/payment/repository"
	"syscall"
	"time"
)

// Run bootstraps every piece of code needed to start the server
func Run(log *zap.SugaredLogger, cnf *config.Config, conn *sql.DB) {
	// Repos
	repo := repository.NewProviderRepo(log, conn)

	// Integrations
	payProvider := intpayment.NewPaymentProvider(log, cnf.ProviderFilePath)
	stores := stores.NewStore(log, cnf.StoresFilePath)

	// Services
	svc := payment.NewPaymentService(log, payProvider, stores, repo)

	// Server setup
	h := v1.NewHandler(log, svc)
	mux := http.NewServeMux()
	logMiddlware := middlwares.LogMiddlware(log)
	headerMiddlware := middlwares.HeaderMiddlware

	mux.HandleFunc("/api/v1/payment/url", headerMiddlware(logMiddlware(h.Payment())))
	svr := http.Server{
		Addr:    cnf.Service.Host + ":" + cnf.Service.Port,
		Handler: mux,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("payment-api http server is closed")
			} else {
				log.Infof("payment-api http server unexpected error: %s", err)
			}
		}
	}()

	// Listening for the stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Errorf("failed to gracefully shutdown the server, error: %v", err)
	}
}
