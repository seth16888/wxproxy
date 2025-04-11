package middleware

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeoutInterceptor 超时拦截器
// 当请求处理时间超过300毫秒时，返回超时错误
// 超时错误的状态码为codes.DeadlineExceeded
// 超时错误的消息为"Deadline exceeded"
func TimeoutInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
		defer cancel()

		ch := make(chan error)
		go func() {
			resp, err = handler(ctxWithTimeout, req)
			ch <- err
		}()

		select {
		case <-ctxWithTimeout.Done():
			cancel()
			err = status.Error(
        codes.DeadlineExceeded,
				fmt.Sprintf("%s: Deadline exceeded", info.FullMethod))
			return
		case <-ch:
		}

		return
	}
}
