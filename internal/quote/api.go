package quote

import (
	"app/internal/catalog"
	"app/internal/order"
	"app/internal/quote/domain"
	"app/internal/quote/handler"
	"app/internal/quote/repository"
	"app/internal/tax"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RouterAPIInitializer() *chi.Mux {
	quoteService := domain.NewQuote(
		repository.NewDynamoQuote(),
		catalog.NewClient(),
		tax.NewClient(),
		order.NewClient(),
	)
	apiHandler := handler.NewAPIHandler(quoteService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/health", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	// Quote Routes
	r.Route("/customers/{customerID}", func(r chi.Router) {
		r.Method("GET", "/quote", handler.BaseHandler(apiHandler.GetQuote()))
		r.Method("POST", "/quote/products", handler.BaseHandler(apiHandler.AddProduct()))
		r.Method("PUT", "/quote/products/{productID}", handler.BaseHandler(apiHandler.UpdateProduct()))
		r.Method("DELETE", "/quote/products/{productID}", handler.BaseHandler(apiHandler.DeleteProduct()))
		r.Method("POST", "/quote", handler.BaseHandler(apiHandler.Process()))
		r.Method("PUT", "/quote", handler.BaseHandler(apiHandler.UpdateAddress()))
		r.Method("PUT", "/quote", handler.BaseHandler(apiHandler.UpdatePayment()))
	})

	return r
}
