package logger

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type Logger struct {
	l *zap.Logger
}

const (
	Key       = "Logger"
	RequestID = "RequestID"
)

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("New zap logger: %w", err)
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(Key).(*Logger); ok {
		return logger
	}
	return nil
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, msg))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, msg))
	}
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, msg))
	}
	l.l.Error(msg, fields...)
}

func LoggerInterceptor(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {

	guid := uuid.New().String()
	ctx = context.WithValue(ctx, RequestID, guid)

	if GetLoggerFromCtx(ctx) == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
		ctx = context.WithValue(ctx, Key, &Logger{logger})
	}

	if logger := GetLoggerFromCtx(ctx); logger != nil {
		logger.Info(ctx, "request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Time("request_time", time.Now()),
		)
	}

	resp, err := handler(ctx, req)

	if logger := GetLoggerFromCtx(ctx); logger != nil {
		if err != nil {
			logger.Error(ctx, "request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		} else {
			logger.Info(ctx, "response",
				zap.String("method", info.FullMethod),
				zap.Any("request", resp),
				zap.Time("request_time", time.Now()),
			)
		}
	}

	return resp, err
}
