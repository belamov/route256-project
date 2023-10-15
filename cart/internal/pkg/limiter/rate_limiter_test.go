package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Wait(t *testing.T) {
	targetRps := 20
	limiter := NewRateLimiter(targetRps)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	increment := atomic.Int32{}
	goroutinesCount := 100
	wg.Add(goroutinesCount)

	ticker := time.NewTicker(time.Second)

	for i := 0; i < goroutinesCount; i++ {
		go func() {
			defer wg.Done()
			err := limiter.Wait(ctx)
			if err != nil {
				return
			}
			increment.Add(1)
		}()
	}

	<-ticker.C
	increment.Store(0)
	<-ticker.C

	result := increment.Load()
	// сложно в тесте получить точное значение, позволим расхождение на +-2
	assert.Greater(t, result, int32(targetRps)-2)
	assert.Less(t, result, int32(targetRps)+2)
	cancel()

	wg.Wait()
}
