package middleware

import (
	"context"

	"github.com/seth16888/wxcommon/helpers"
	"github.com/seth16888/wxproxy/internal/consts"
	"google.golang.org/grpc"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// 生成Request ID
		requestID := helpers.UUID()

		ctx = context.WithValue(ctx, consts.RequestIdKey, requestID)

		// 继续处理请求
		return handler(ctx, req)
	}
}
