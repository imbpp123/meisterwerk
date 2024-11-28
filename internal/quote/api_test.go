package quote_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"app/internal/catalog"
	"app/internal/customer"
	"app/internal/order"
	"app/internal/quote/domain"
	"app/internal/quote/handler"
	"app/internal/quote/repository"
	"app/internal/tax"
)

type testApiHandle struct {
	handler         *handler.APIHandler
	customerService *customer.Client
}

func newTestApiHandler() *testApiHandle {
	quoteService := domain.NewQuote(
		repository.NewDynamoQuote(),
		catalog.NewClient(),
		tax.NewClient(),
		order.NewClient(),
	)

	return &testApiHandle{
		handler: handler.NewAPIHandler(quoteService),
	}
}

// api test example...
// check contract here, also it can be used as integration test...
func TestApiHandlerGetQuoteNewQuote(t *testing.T) {
	// arrange
	customerUUID := uuid.NewString()
	tc := newTestApiHandler()

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/customers/%s/quote", customerUUID), nil)
	assert.NoError(t, err)

	r := chi.NewRouter()
	r.Use(handler.CustomerCtxMiddleware(tc.customerService))
	r.Method("GET", "/customers/{customerID}/quote", handler.BaseHandler(tc.handler.GetQuote()))

	fmt.Println(r.Routes())

	// act
	r.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	quote := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &quote))

	expectedQuote := map[string]interface{}{
		"amount":     0,
		"tax_amount": 0,
	}

	assert.Equal(t, expectedQuote, quote)
}
