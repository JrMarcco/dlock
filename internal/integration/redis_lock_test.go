package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/JrMarcco/dlock/redis"
	"github.com/JrMarcco/easy-kit/retry"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRedisDlock(t *testing.T) {
	rc := redis.NewClient(&redis.Options{
		Addr:     "192.168.3.3:6379",
		Password: "<passwd>",
		DB:       0,
	})

	err := rc.Ping(context.Background()).Err()
	require.NoError(t, err)

	strategy, _ := retry.NewExponentialBackoffStrategy(10*time.Millisecond, 50*time.Millisecond, 8)

	suite.Run(t, &RedisDlockTestSuite{
		DlockTestSuite: &DlockTestSuite{
			dc: rdlock.NewDClientBuilder(rc).WithRetryStrategy(strategy).Build(),
		},
		rc: rc,
	})
}

type RedisDlockTestSuite struct {
	*DlockTestSuite
	rc redis.Cmdable
}

func (s *RedisDlockTestSuite) SetupSuite() {}

func (s *RedisDlockTestSuite) TearDownSuite() {
	ctx := context.Background()
	keys, err := s.rc.Keys(ctx, fmt.Sprintf("%s*", keyPrefix)).Result()
	require.NoError(s.T(), err)

	if len(keys) > 0 {
		err = s.rc.Del(ctx, keys...).Err()
		require.NoError(s.T(), err)
	}
}
