package customer

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) IsActive(ctx context.Context, ID uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}
