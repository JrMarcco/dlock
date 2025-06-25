package redis

import (
	"context"
	"time"

	"github.com/JrMarcco/dlock"
	"github.com/JrMarcco/easy-kit/bean/option"
	"github.com/JrMarcco/easy-kit/retry"
	"github.com/redis/go-redis/v9"
)

type ClientBuilder struct {
	rc   redis.Cmdable
	opts []option.Opt[Lock]
}

func (b *ClientBuilder) Build() *Client {
	return &Client{
		rc:   b.rc,
		opts: b.opts,
	}
}

// WithRetryStrategy 指定重试策略
func (b *ClientBuilder) WithRetryStrategy(strategy retry.Strategy) *ClientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Lock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Lock) {
		lock.retryStrategy = strategy
	})
	return b
}

// WithValFunc 指定值生成方法
func (b *ClientBuilder) WithValFunc(valFunc func() string) *ClientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Lock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Lock) {
		lock.valFunc = valFunc
	})
	return b
}

func (b *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
	if b.opts == nil {
		b.opts = make([]option.Opt[Lock], 0, 2)
	}
	b.opts = append(b.opts, func(lock *Lock) {
		lock.timeout = timeout
	})
	return b
}

func NewClientBuilder(rc redis.Cmdable) *ClientBuilder {
	return &ClientBuilder{
		rc: rc,
	}
}

type Client struct {
	rc   redis.Cmdable
	opts []option.Opt[Lock]
}

func (c *Client) NewLock(_ context.Context, key string, expiration time.Duration) (dlock.Lock, error) {
	return NewLock(c.rc, key, expiration, c.opts...)
}
