package order

import (
	"context"
	"errors"

	"app/internal/quote/types"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Process(ctx context.Context, quote *types.Quote) error {
	return errors.New("not implemeted")
}
