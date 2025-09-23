package util

import (
	"context"
	"time"
)

const DefaultDBTimeout = 5 * time.Second
const DefaultRedisTimeout = 2 * time.Second

// NewDefaultDBContext creates a new context with the default database timeout.
func NewDefaultDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultDBTimeout)
}

// NewDBContextWith creates a new context with a custom timeout duration
func NewDBContextWith(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// NewDefaultRedisContext creates a new context with the default Redis timeout.
func NewDefaultRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultRedisTimeout)
}

// NewRedisContextWith creates a new context with a custom timeout duration
func NewRedisContextWith(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
