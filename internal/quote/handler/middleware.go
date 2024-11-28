package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type (
	customerIDCtx struct{}

	customerService interface {
		IsActive(ctx context.Context, customerUUID uuid.UUID) (bool, error)
	}
)

var (
	errCustomerDisabled = errors.New("customer is disabled")
)

func CustomerCtxMiddleware(customerService customerService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			customerID, err := getParamUUID(r, URLCustomerIDParameter)
			if err != nil {
				respondError(w, err)
				return
			}

			active, err := customerService.IsActive(r.Context(), customerID)
			if err != nil {
				respondError(w, err)
				return
			}
			if !active {
				respondError(w, errCustomerDisabled)
				return
			}

			ctx := context.WithValue(r.Context(), customerIDCtx{}, customerID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
