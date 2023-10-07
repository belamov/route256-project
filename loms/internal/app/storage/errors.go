package storage

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInsufficientStocks = errors.New("not enough stocks")
)
