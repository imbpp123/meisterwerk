package types

import (
	"time"

	"github.com/google/uuid"
)

type QuoteStatus string

const (
	QuoteStatusDraft      QuoteStatus = "draft"
	QuoteStatusProcessing QuoteStatus = "processing"
	QuoteStatusDone       QuoteStatus = "done"
)

type Quote struct {
	UUID        uuid.UUID
	CustomerID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      QuoteStatus
	Amount      float64
	TaxAmount   float64
	TotalAmount float64
	Address     *Address
	Payment     *Payment
	Products    []Product
}

type Address struct {
	Address string
	City    string
	Country string
}

type Payment struct {
	PaymentMethod string
}

type Product struct {
	ProductID   uuid.UUID
	Quantity    int
	Amount      float64
	TaxAmount   float64
	TotalAmount float64
}

type ProductAdd struct {
	ProductID uuid.UUID
	Quantity  int
}

type ProductUpdate struct {
	Quantity int
}

func NewQuote(
	UUID uuid.UUID,
	customerID uuid.UUID,
) *Quote {
	return &Quote{
		UUID:        UUID,
		CustomerID:  customerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      QuoteStatusDraft,
		Amount:      0,
		TaxAmount:   0,
		TotalAmount: 0,
		Address:     nil,
		Payment:     nil,
		Products:    nil,
	}
}
