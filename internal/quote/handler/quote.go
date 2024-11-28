package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"app/internal/quote/types"
)

const (
	URLCustomerIDParameter string = "customerID"
	URLProductIDParameter  string = "productID"
)

type (
	quoteResponse struct {
		ID          uuid.UUID         `json:"id"`
		Address     addressResponse   `json:"address"`
		Payment     paymentResponse   `json:"payment"`
		Products    []productResponse `json:"products"`
		Amount      float64           `json:"amount"`
		TaxAmount   float64           `json:"tax_amount"`
		TotalAmount float64           `json:"total_amount"`
	}

	addressResponse struct {
		Address string `json:"address"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	paymentResponse struct {
		PaymentMethod string `json:"payment_method"`
	}

	productResponse struct {
		ID          uuid.UUID `json:"product_id"`
		Quantity    int       `json:"qty"`
		Amount      float64   `json:"amount"`
		TaxAmount   float64   `json:"tax_amount"`
		TotalAmount float64   `json:"total_amount"`
	}

	addressRequest struct {
		Address string `json:"address"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	paymentRequest struct {
		PaymentMethod string `json:"payment_method"`
	}

	productAddRequest struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"qty"`
	}

	productUpdateRequest struct {
		Quantity int `json:"qty"`
	}
)

type (
	quoteService interface {
		AddProduct(ctx context.Context, customerUUID uuid.UUID, product *types.ProductAdd) error
		UpdateProduct(ctx context.Context, customerUUID uuid.UUID, productUUID uuid.UUID, product *types.ProductUpdate) error
		LoadDraftByCustomer(ctx context.Context, customerUUID uuid.UUID) (*types.Quote, error)
		ProcessByCustomerID(ctx context.Context, customerUUID uuid.UUID) error
		RemoveProduct(ctx context.Context, customerUUID uuid.UUID, productID uuid.UUID) error
		SaveAddress(ctx context.Context, customerUUID uuid.UUID, address *types.Address) error
		SavePayment(ctx context.Context, customerUUID uuid.UUID, payment *types.Payment) error
	}

	APIHandler struct {
		quoteService quoteService
	}
)

func NewAPIHandler(quoteService quoteService) *APIHandler {
	return &APIHandler{
		quoteService: quoteService,
	}
}

func (q *APIHandler) GetQuote() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::GetQuote : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) UpdateAddress() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateAddress : %w", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateAddress : %w: %w", errBodyRead, err)
		}

		var request addressRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateAddress : %w: %w", errBodyRead, err)
		}

		err = q.quoteService.SaveAddress(r.Context(), customerID, &types.Address{
			Address: request.Address,
			City:    request.City,
			Country: request.Country,
		})
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateAddress : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) UpdatePayment() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdatePayment : %w", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdatePayment : %w: %w", errBodyRead, err)
		}

		var request paymentRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdatePayment : %w: %w", errBodyRead, err)
		}

		err = q.quoteService.SavePayment(r.Context(), customerID, &types.Payment{
			PaymentMethod: request.PaymentMethod,
		})
		if err != nil {
			return fmt.Errorf("APIHandler::UpdatePayment : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) Process() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::Process : %w", err)
		}

		if err := q.quoteService.ProcessByCustomerID(r.Context(), customerID); err != nil {
			return fmt.Errorf("APIHandler::Process : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) AddProduct() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::AddProduct : %w", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("APIHandler::AddProduct : %w: %w", errBodyRead, err)
		}

		var request productAddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			return fmt.Errorf("APIHandler::AddProduct : %w: %w", errBodyRead, err)
		}

		err = q.quoteService.AddProduct(r.Context(), customerID, &types.ProductAdd{})
		if err != nil {
			return fmt.Errorf("APIHandler::AddProduct : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) UpdateProduct() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateProduct : %w", err)
		}

		productID, err := getParamUUID(r, URLProductIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateProduct : %w", err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateProduct : %w: %w", errBodyRead, err)
		}

		var request productUpdateRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			return fmt.Errorf("APIHandler::UpdateProduct : %w: %w", errBodyRead, err)
		}

		err = q.quoteService.UpdateProduct(r.Context(), customerID, productID, &types.ProductUpdate{
			Quantity: request.Quantity,
		})
		if err != nil {
			return fmt.Errorf("APIHandler::DeleteProduct : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) DeleteProduct() BaseHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		customerID, err := getParamUUID(r, URLCustomerIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::DeleteProduct : %w", err)
		}

		productID, err := getParamUUID(r, URLProductIDParameter)
		if err != nil {
			return fmt.Errorf("APIHandler::DeleteProduct : %w", err)
		}

		if err := q.quoteService.RemoveProduct(r.Context(), customerID, productID); err != nil {
			return fmt.Errorf("APIHandler::DeleteProduct : %w", err)
		}

		return q.respondQuote(r.Context(), w, customerID)
	}
}

func (q *APIHandler) respondQuote(ctx context.Context, w http.ResponseWriter, customerID uuid.UUID) error {
	quote, err := q.quoteService.LoadDraftByCustomer(ctx, customerID)
	if err != nil {
		return fmt.Errorf("APIHandler::respondQuote : %w", err)
	}

	return respond(
		w,
		quote, // need to change it here...
		http.StatusOK,
	)
}
