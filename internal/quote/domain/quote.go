package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"app/internal/catalog"
	"app/internal/quote/types"
)

type (
	orderClient interface {
		Process(ctx context.Context, quote *types.Quote) error
	}

	catalogClient interface {
		GetProductByID(ctx context.Context, productID uuid.UUID) (*catalog.Product, error)
	}

	taxClient interface {
		CalculateTaxes(ctx context.Context, taxRateID string, amount float64) (float64, error)
	}

	quoteRepository interface {
		FindByCustomerAndStatus(ctx context.Context, customerUUID uuid.UUID, status types.QuoteStatus) (*types.Quote, error)
		Save(ctx context.Context, quote *types.Quote) error
	}

	Quote struct {
		repository quoteRepository
		catalog    catalogClient
		taxes      taxClient
		order      orderClient
	}
)

func NewQuote(
	repository quoteRepository,
	catalog catalogClient,
	taxes taxClient,
	order orderClient,
) *Quote {
	return &Quote{
		repository: repository,
		catalog:    catalog,
		taxes:      taxes,
		order:      order,
	}
}

// AddProduct adds a new product to the customer's draft quote.
// If the quote doesn't exist, it creates a new one.
func (q *Quote) AddProduct(ctx context.Context, customerUUID uuid.UUID, product *types.ProductAdd) error {
	return q.withDraft(ctx, customerUUID, func(quote *types.Quote) error {
		quote.Products = append(quote.Products, types.Product{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})

		if err := q.refresh(ctx, quote); err != nil {
			return fmt.Errorf("Domain::Quote::AddProduct : %w", err)
		}

		return nil
	})
}

// UpdateProduct updates the quantity of a product in the customer's draft quote.
// Returns an error if the product is not found in the quote.
func (q *Quote) UpdateProduct(ctx context.Context, customerUUID uuid.UUID, productUUID uuid.UUID, product *types.ProductUpdate) error {
	return q.withDraft(ctx, customerUUID, func(quote *types.Quote) error {
		isFound := false
		n := len(quote.Products)
		for i := 0; i < n; i++ {
			if quote.Products[i].ProductID == productUUID {
				quote.Products[i].Quantity = product.Quantity
				isFound = true
				break
			}
		}
		if !isFound {
			return types.ErrQuoteProductNotFound
		}

		if err := q.refresh(ctx, quote); err != nil {
			return fmt.Errorf("Domain::Quote::UpdateProduct : %w", err)
		}

		return nil
	})
}

// LoadDraftByCustomer retrieves the customer draft quote or creates a new one if it doesn't exist.
func (q *Quote) LoadDraftByCustomer(ctx context.Context, customerUUID uuid.UUID) (*types.Quote, error) {
	quote, err := q.repository.FindByCustomerAndStatus(ctx, customerUUID, types.QuoteStatusDraft)
	if err != nil {
		if errors.Is(err, types.ErrQuoteNotFound) {
			return types.NewQuote(uuid.New(), customerUUID), nil
		}

		return nil, fmt.Errorf("Domain::Quote::LoadDraftByCustomer : %w", err)
	}

	return quote, nil
}

// ProcessByCustomerID processes the customer draft quote and marks it as done.
func (q *Quote) ProcessByCustomerID(ctx context.Context, customerUUID uuid.UUID) error {
	quote, err := q.LoadDraftByCustomer(ctx, customerUUID)
	if err != nil {
		return fmt.Errorf("Domain::Quote::ProcessByCustomerID : %w", err)
	}

	quote.Status = types.QuoteStatusProcessing
	quote.UpdatedAt = time.Now()

	if err := q.repository.Save(ctx, quote); err != nil {
		return fmt.Errorf("Domain::Quote::ProcessByCustomerID : %w", err)
	}

	if err := q.order.Process(ctx, quote); err != nil {
		return fmt.Errorf("Domain::Quote::ProcessByCustomerID : %w", err)
	}

	// it can be changed via events from Order Scheduling Service
	quote.Status = types.QuoteStatusDone
	quote.UpdatedAt = time.Now()

	if err := q.repository.Save(ctx, quote); err != nil {
		return fmt.Errorf("Domain::Quote::ProcessByCustomerID : %w", err)
	}

	return nil
}

// RemoveProduct removes a product from the customer's draft quote.
// Returns an ErrQuoteProductNotFound error if the product is not found.
func (q *Quote) RemoveProduct(ctx context.Context, customerUUID uuid.UUID, productID uuid.UUID) error {
	return q.withDraft(ctx, customerUUID, func(quote *types.Quote) error {
		filtered := make([]types.Product, 0)
		for _, product := range quote.Products {
			if product.ProductID != productID {
				filtered = append(filtered, product)
			}
		}
		if len(filtered) == len(quote.Products) {
			return types.ErrQuoteProductNotFound
		}
		quote.Products = filtered

		if err := q.refresh(ctx, quote); err != nil {
			return fmt.Errorf("Domain::Quote::RemoveProduct : %w", err)
		}

		return nil
	})
}

// SaveAddress saves the customer's address in the draft quote.
func (q *Quote) SaveAddress(ctx context.Context, customerUUID uuid.UUID, address *types.Address) error {
	return q.withDraft(ctx, customerUUID, func(quote *types.Quote) error {
		quote.Address = address

		return nil
	})
}

// SavePayment saves the customer's payment details in the draft quote.
func (q *Quote) SavePayment(ctx context.Context, customerUUID uuid.UUID, payment *types.Payment) error {
	return q.withDraft(ctx, customerUUID, func(quote *types.Quote) error {
		quote.Payment = payment

		return nil
	})
}

// calculateProduct calculates the tax and total amount for a product in the quote.
func (q *Quote) calculateProduct(ctx context.Context, product *types.Product) error {
	productInfo, err := q.catalog.GetProductByID(ctx, product.ProductID)
	if err != nil {
		return fmt.Errorf("Domain::Quote::calculateProduct : %w", err)
	}

	product.Amount = float64(product.Quantity) * productInfo.Price
	product.TaxAmount, err = q.taxes.CalculateTaxes(ctx, productInfo.TaxRateID, product.Amount)
	if err != nil {
		return fmt.Errorf("Domain::Quote::calculateProduct : %w", err)
	}

	product.TotalAmount = product.Amount + product.TaxAmount
	return nil
}

// refresh recalculates the totals for the quote based on its products.
func (q *Quote) refresh(ctx context.Context, quote *types.Quote) error {
	quote.Amount, quote.TaxAmount, quote.TotalAmount = 0, 0, 0

	for i := range quote.Products {
		if err := q.calculateProduct(ctx, &quote.Products[i]); err != nil {
			return fmt.Errorf("Domain::Quote::refresh : %w", err)
		}

		quote.Amount += quote.Products[i].Amount
		quote.TaxAmount += quote.Products[i].TaxAmount
		quote.TotalAmount += quote.Products[i].TotalAmount
	}

	return nil
}

// withDraft executes an action on the customer's draft quote and saves the updated quote.
func (q *Quote) withDraft(ctx context.Context, customerUUID uuid.UUID, action func(*types.Quote) error) error {
	return q.withLock(ctx, customerUUID, func() error {
		quote, err := q.LoadDraftByCustomer(ctx, customerUUID)
		if err != nil {
			return fmt.Errorf("Domain::Quote::withDraft : %w", err)
		}

		if err := action(quote); err != nil {
			return fmt.Errorf("Domain::Quote::withDraft : %w", err)
		}

		quote.UpdatedAt = time.Now()
		if err := q.repository.Save(ctx, quote); err != nil {
			return fmt.Errorf("Domain::Quote::withDraft : %w", err)
		}
		return nil
	})
}

func (q *Quote) withLock(ctx context.Context, customerUUID uuid.UUID, action func() error) error {
	// Implementation of distributed lock should go here using AWS Memcached-Redis-similar product
	return action()
}
