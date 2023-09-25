package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"payment-api/internal/services/payment"
)

type Payment interface {
	PaymentUrl(ctx context.Context, providerID string) (string, error)
	StoresUrls(ctx context.Context) ([]map[string]string, error)
}

type Handler struct {
	log        *zap.SugaredLogger
	paymentSvc Payment
}

func NewHandler(log *zap.SugaredLogger, paymentSvc Payment) *Handler {
	return &Handler{log: log, paymentSvc: paymentSvc}
}

// Payment endpoint for retrieving url for the provided productID
func (h *Handler) Payment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prodID := r.URL.Query().Get("productID")
		if prodID == "" {
			h.log.Errorf("failed to retrieve 'productID' parameter")
			writeJson(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "Missing productID parameter"})
			return
		}
		url, err := h.paymentSvc.PaymentUrl(r.Context(), prodID)

		if err != nil {
			h.log.Errorf("failed to receive payment url")
			switch {
			case errors.Is(err, payment.ErrUuidInvalidFormat), errors.Is(err, payment.ErrNotFound):
				writeJson(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "Provided parameter has bad format"})
			case errors.Is(err, payment.ErrProvider):
				urls, err := h.paymentSvc.StoresUrls(r.Context())
				if err != nil {
					writeJson(w, http.StatusBadRequest, map[string]any{"code": http.StatusInternalServerError, "message": "Oops, something went wrong"})
					return
				}
				writeJson(w, http.StatusBadRequest, map[string]any{"code": http.StatusOK, "stores_urls": urls})
			}
			return
		}
		writeJson(w, http.StatusBadRequest, map[string]any{"code": http.StatusOK, "data": url})
	}
}

// writeJson peforms write to the response writer and header as json content-type
func writeJson(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	enc, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, _ = w.Write(enc)
}
