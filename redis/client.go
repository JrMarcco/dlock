package rdlock

import (
	"context"
	"time"

	"github.com/JrMarcco/dlock"
	"github.com/JrMarcco/easy-kit/bean/option"
	"github.com/JrMarcco/easy-kit/retry"
	"github.com/redis/go-redis/v9"
)

type DclientBuilder struct {
	rc   redis.Cmdable
	opts []option.Opt[Dlock]
}

func (b *DclientBuilder) Build() *Dclient {
	return &Dclient{
		rc:   b.rc,
		opts: b.opts,
	}
}

// WithRetryStrategy 指定重试策略
func (b *DclientBuilder) WithRetryStrategy(strategy retry.Strategy) *DclientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Dlock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Dlock) {
		lock.retryStrategy = strategy
	})
	return b
}

// WithValFunc 指定值生成方法
func (b *DclientBuilder) WithValFunc(valFunc func() string) *DclientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Dlock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Dlock) {
		lock.valFunc = valFunc
	})
	return b
}

func (b *DclientBuilder) WithTimeout(timeout time.Duration) *DclientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Dlock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Dlock) {
		lock.timeout = timeout
	})
	return b
}

func NewDClientBuilder(rc redis.Cmdable) *DclientBuilder {
	return &DclientBuilder{
		rc: rc,
	}
}

var _ dlock.Dclient = (*Dclient)(nil)

type Dclient struct {
	rc   redis.Cmdable
	opts []option.Opt[Dlock]
}

func (c *Dclient) NewDlock(_ context.Context, key string, expiration time.Duration) (dlock.Dlock, error) {
	return NewDlock(c.rc, key, expiration, c.opts...)
}
