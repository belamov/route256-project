package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type Limiter interface {
	Wait(ctx context.Context) error
}

func RateLimitClientInterceptor(limiter Limiter) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := limiter.Wait(ctx)
		if err != nil {
			return fmt.Errorf("cant wait for limiter to allow request: %w", err)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
