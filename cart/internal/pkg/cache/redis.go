package cache

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"route256/cart/internal/app/models"
)

type Redis struct {
	cache *cache.Cache
}

func NewRedis(ctx context.Context, wg *sync.WaitGroup, shards []string) *Redis {
	addrs := make(map[string]string)
	for i, shard := range shards {
		addrs["shard"+strconv.Itoa(i)] = shard
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
	})

	mycache := cache.New(&cache.Options{
		Redis: ring,
	})

	go func() {
		defer wg.Done()
		<-ctx.Done()
		err := ring.Close()
		if err != nil {
			log.Error().Err(err).Msg("cant close redis connections")
		}
	}()

	return &Redis{cache: mycache}
}

func (r Redis) GetCartItems(ctx context.Context, userId int64) ([]models.CartItemWithInfo, error) {
	items := make([]models.CartItemWithInfo, 0)
	err := r.cache.Get(ctx, strconv.FormatInt(userId, 10), &items)
	if errors.Is(err, cache.ErrCacheMiss) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r Redis) SetCartItems(ctx context.Context, userId int64, items []models.CartItemWithInfo) error {
	cacheItem := &cache.Item{
		Key:   r.keyFromUserId(userId),
		Value: items,
	}
	return r.cache.Set(cacheItem)
}

func (r Redis) Invalidate(ctx context.Context, userId int64) error {
	return r.cache.Delete(ctx, r.keyFromUserId(userId))
}

func (r Redis) keyFromUserId(userId int64) string {
	return strconv.FormatInt(userId, 10)
}
