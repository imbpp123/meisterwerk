package repository

import (
	"app/internal/quote/types"
	"context"
	"errors"

	"github.com/google/uuid"
)

type DynamoQuote struct {
}

func NewDynamoQuote() *DynamoQuote {
	return &DynamoQuote{}
}

func (d *DynamoQuote) FindByCustomerAndStatus(ctx context.Context, customerUUID uuid.UUID, status types.QuoteStatus) (*types.Quote, error) {
	return nil, errors.New("not implemented")
}

func (d *DynamoQuote) Save(ctx context.Context, quote *types.Quote) error {
	return errors.New("not implemented")
}
