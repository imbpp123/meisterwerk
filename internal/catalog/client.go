package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type (
	Product struct {
		ProductID uuid.UUID
		Price     float64
		TaxRateID string
	}

	Client struct {
	}
)

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetProductByID(ctx context.Context, productID uuid.UUID) (*Product, error) {
	return nil, errors.New("not implemented")
}
