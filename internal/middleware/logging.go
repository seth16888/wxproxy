package middleware

import (
	"context"
	"time"

	"github.com/seth16888/wxproxy/internal/consts"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, // 上下文对象
		req any, // 请求对象
		info *grpc.UnaryServerInfo, // 服务信息
		handler grpc.UnaryHandler, // 处理函数
	) (resp any, err error) {
		start := time.Now().UnixMilli()

		resp, err = handler(ctx, req)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Int64("latency", time.Now().UnixMilli()-start),
			zap.Error(err),
		}
		// RequestID
		requestID, ok := ctx.Value(consts.RequestIdKey).(string)
		if ok {
			fields = append(fields, zap.String("requestId", requestID))
		}

		log.Info("request", fields...)

		return resp, err
	}
}
