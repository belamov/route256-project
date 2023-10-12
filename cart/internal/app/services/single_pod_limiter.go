package services

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type SinglePodLimiter struct {
	limiter *rate.Limiter
}

func NewSinglePodLimiter(targetRps int) *SinglePodLimiter {
	return &SinglePodLimiter{
		limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(targetRps)), 1),
	}
}

func (s *SinglePodLimiter) Wait(ctx context.Context) error {
	return s.limiter.Wait(ctx)
}
