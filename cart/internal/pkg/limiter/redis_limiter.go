package limiter

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisRateLimiter struct {
	limiter   *redis_rate.Limiter
	targetRps int
	key       string
}

func NewRedisRateLimiter(ctx context.Context, wg *sync.WaitGroup, redisAddr string, targetRps int) *RedisRateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	log.Info().Msg("connected to redis on " + redisAddr)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info().Msg("closing connection to redis...")
		err := rdb.Close()
		if err != nil {
			log.Err(err).Msg("failed to close redis connection: ")
			return
		}
		log.Info().Msg("closed connection to redis")
	}()

	return &RedisRateLimiter{
		limiter:   redis_rate.NewLimiter(rdb),
		targetRps: targetRps,
		key:       "product-client-rate",
	}
}

func (r RedisRateLimiter) Wait(ctx context.Context) error {
	timer := time.NewTimer(0)
	for {
		select {
		case <-timer.C:
			timer.Stop()
			res, err := r.tryAcquireToken(ctx)
			if err != nil {
				return err
			}
			if res.Allowed > 0 {
				return nil
			}
			timer.Reset(res.RetryAfter)
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
}

func (r RedisRateLimiter) tryAcquireToken(ctx context.Context) (*redis_rate.Result, error) {
	return r.limiter.Allow(ctx, r.key, redis_rate.PerSecond(r.targetRps))
}
