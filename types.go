package dlock

import (
	"context"
	"time"
)

type Dlock interface {
	TryLock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Refresh(ctx context.Context) error
}

type Dclient interface {
	NewDlock(ctx context.Context, key string, expiration time.Duration) (Dlock, error)
}
