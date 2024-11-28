package handler

import (
	"app/internal/quote/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type (
	BaseHandler func(w http.ResponseWriter, r *http.Request) error

	errorResponse struct {
		Status  int    `json:"-"`
		Message string `json:"message"`
	}
)

var (
	errorResponseCodes = map[error]errorResponse{
		errMissedRequiredParameter: {
			Status:  http.StatusBadRequest,
			Message: "missed required parameter",
		},
		errInvalidParameter: {
			Status:  http.StatusBadRequest,
			Message: "parameter is invalid",
		},
		errBodyRead: {
			Status:  http.StatusBadRequest,
			Message: "parameter is invalid",
		},
		types.ErrQuoteNotFound: {
			Status:  http.StatusNotFound,
			Message: "quote not found",
		},
	}

	errMissedRequiredParameter = errors.New("missing required parameter")
	errInvalidParameter        = errors.New("invalid parameter")
	errBodyRead                = errors.New("invalid body")
)

func (fn BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		respondError(w, err)
	}
}

func respondError(w http.ResponseWriter, err error) error {
	for e, apiError := range errorResponseCodes {
		if errors.Is(err, e) {
			return respond(w, apiError, apiError.Status)
		}
	}

	return respond(w, nil, http.StatusInternalServerError)
}

// Respond writes the given data to an HTTP response with a status code.
func respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.WriteHeader(statusCode)

	if data == nil || statusCode == http.StatusNoContent {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

func getParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	parameterIDStr := chi.URLParam(r, name)
	if parameterIDStr == "" {
		return uuid.UUID{}, errMissedRequiredParameter
	}

	parameterUUID, err := uuid.Parse(parameterIDStr)
	if err != nil {
		return uuid.UUID{}, errInvalidParameter
	}

	return parameterUUID, nil
}
