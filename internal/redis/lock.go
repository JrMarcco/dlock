package redis

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/JrMarcco/dlock"
	"github.com/JrMarcco/easy-kit/bean/option"
	"github.com/JrMarcco/easy-kit/retry"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/lock.lua
	lockLua string
	//go:embed lua/unlock.lua
	unLockLua string
	//go:embed lua/refresh.lua
	refreshLua string
)

var _ dlock.Lock = (*Lock)(nil)

type Lock struct {
	client redis.Cmdable

	key     string
	val     string
	valFunc func() string // 生成 val 的方法，默认 uuid

	timeout       time.Duration  // 单词加锁超时时间，默认 100ms
	expiration    time.Duration  // 过期时间
	retryStrategy retry.Strategy // 加锁重试策略
}

// TryLock 尝试获取分布式锁，当失败时候会根据重试策略进行重试。
// 默认重试策略为指数退避策略（初始间隔 100ms，最大间隔 1s，最大重试次数 8 次）。
func (l *Lock) TryLock(ctx context.Context) error {
	return retry.Retry(ctx, l.retryStrategy, func() error {
		lockCtx, cancel := context.WithTimeout(ctx, l.timeout)
		defer cancel()

		res, err := l.client.Eval(
			lockCtx,
			lockLua,
			[]string{l.key},
			l.val,
			l.expiration.Seconds(),
		).Bool()
		if err != nil {
			return err
		}

		if res {
			// 加锁成功
			return nil
		}
		return dlock.ErrLockIsHeld
	})
}

// Unlock 释放锁，失败不会重试。
func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, unLockLua, []string{l.key}, l.val).Bool()
	if errors.Is(err, redis.Nil) {
		// key 不存在
		return dlock.ErrReleaseNotHeld
	}
	if err != nil {
		return err
	}

	if res {
		return nil
	}
	return dlock.ErrReleaseNotHeld
}

// Refresh 刷新锁的过期时间，失败不会重试。
func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, refreshLua, []string{l.key}, l.val, l.expiration.Seconds()).Bool()
	if err != nil {
		return err
	}
	if res {
		return nil
	}
	return dlock.ErrRefreshNotHeld
}

func NewLock(client redis.Cmdable, key string, expiration time.Duration, opts ...option.Opt[Lock]) (*Lock, error) {
	// 默认指数退避重试策略
	strategy, _ := retry.NewExponentialBackoffStrategy(100*time.Millisecond, time.Second, 8)

	lock := &Lock{
		client: client,
		key:    key,
		valFunc: func() string {
			return uuid.New().String()
		},
		timeout:       100 * time.Millisecond,
		expiration:    expiration,
		retryStrategy: strategy,
	}

	for _, opt := range opts {
		opt(lock)
	}

	lock.val = lock.valFunc()
	return lock, nil
}
