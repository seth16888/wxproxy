package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ClientDisconnectInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {

		ch := make(chan error)

		go func() {
			resp, err = handler(ctx, req)
			ch <- err
		}()

		select {
		case <-ctx.Done(): // 客户端断开连接
			err = status.Error(codes.Canceled, fmt.Sprintf("%s: Request canceled", info.FullMethod))
			return
		case <-ch:
		}

		return
	}
}
