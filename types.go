package dlock

import (
	"context"
	"time"
)

type Lock interface {
	TryLock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Refresh(ctx context.Context) error
}

type Client interface {
	New(ctx context.Context, key string, expiration time.Duration) (Lock, error)
}
