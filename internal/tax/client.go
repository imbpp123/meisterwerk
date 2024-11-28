package tax

import (
	"context"
	"errors"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (t *Client) CalculateTaxes(ctx context.Context, taxRateID string, amount float64) (float64, error) {
	return 0, errors.New("not implemented")
}
