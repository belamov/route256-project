package services

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(targetRps int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(targetRps)), 1),
	}
}

func (s *RateLimiter) Wait(ctx context.Context) error {
	return s.limiter.Wait(ctx)
}
