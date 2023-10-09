package services

import "context"

type MockTransactor struct{}

func (m MockTransactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	return tFunc(ctx)
}
