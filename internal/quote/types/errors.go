package types

import "errors"

var (
	ErrQuoteNotFound        = errors.New("quote not found")
	ErrQuoteProductNotFound = errors.New("quote product not found")
	ErrQuoteUnchangeable    = errors.New("quote can not be changed")
)
