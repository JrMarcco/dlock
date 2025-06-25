package integration

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JrMarcco/dlock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const keyPrefix = "test:dlock"

type DlockTestSuite struct {
	suite.Suite
	dc dlock.Dclient
}

func (s *DlockTestSuite) TestTryLock() {
	tcs := []struct {
		name    string
		dlFunc  func(t *testing.T) dlock.Dlock
		before  func(t *testing.T)
		wantErr error
	}{
		{
			name: "basic",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:lock:basic", keyPrefix), time.Minute)
				require.NoError(s.T(), err)
				return dl
			},
			before:  func(t *testing.T) {},
			wantErr: nil,
		}, {
			name: "held by another process",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:lock:held", keyPrefix), time.Minute)
				require.NoError(s.T(), err)
				return dl
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:lock:held", keyPrefix), time.Minute)
				require.NoError(s.T(), err)

				lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
				defer lockCancel()
				err = dl.TryLock(lockCtx)
				require.NoError(s.T(), err)
			},
			wantErr: dlock.ErrLockIsHeld,
		}, {
			name: "succeed to acquire lock after lock holder crash",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:lock:succeed", keyPrefix), time.Minute)
				require.NoError(s.T(), err)
				return dl
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:lock:succeed", keyPrefix), time.Second)
				require.NoError(s.T(), err)

				err = dl.TryLock(ctx)
				require.NoError(s.T(), err)

				time.Sleep(time.Second)
			},
			wantErr: nil,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)

			dl := tc.dlFunc(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := dl.TryLock(ctx)
			cancel()
			assert.True(t, errors.Is(err, tc.wantErr))
		})
	}
}

func (s *DlockTestSuite) TestUnLock() {
	tcs := []struct {
		name    string
		dlFunc  func(t *testing.T) dlock.Dlock
		before  func(t *testing.T)
		wantErr error
	}{
		{
			name: "basic",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:unlock:basic", keyPrefix), time.Minute)
				require.NoError(s.T(), err)

				lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
				defer lockCancel()
				err = dl.TryLock(lockCtx)
				require.NoError(s.T(), err)
				return dl
			},
			before:  func(t *testing.T) {},
			wantErr: nil,
		}, {
			name: "release not held",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:unlock:nh", keyPrefix), time.Minute)
				require.NoError(s.T(), err)
				return dl
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:unlock:nh", keyPrefix), time.Minute)
				require.NoError(s.T(), err)

				lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
				defer lockCancel()
				err = dl.TryLock(lockCtx)
				require.NoError(s.T(), err)
			},
			wantErr: dlock.ErrReleaseNotHeld,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)

			dl := tc.dlFunc(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := dl.Unlock(ctx)
			cancel()

			assert.True(t, errors.Is(err, tc.wantErr))
		})
	}
}

func (s *DlockTestSuite) TestRefresh() {
	tcs := []struct {
		name    string
		dlFunc  func(t *testing.T) dlock.Dlock
		before  func(t *testing.T)
		wantErr error
	}{
		{
			name: "basic",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:refresh:basic", keyPrefix), time.Minute)
				require.NoError(s.T(), err)

				lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
				defer lockCancel()
				err = dl.TryLock(lockCtx)
				require.NoError(s.T(), err)
				return dl
			},
			before:  func(t *testing.T) {},
			wantErr: nil,
		}, {
			name: "refresh not held",
			dlFunc: func(t *testing.T) dlock.Dlock {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:refresh:nh", keyPrefix), time.Minute)
				require.NoError(s.T(), err)
				return dl
			},
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:refresh:nh", keyPrefix), time.Minute)
				require.NoError(s.T(), err)

				lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
				defer lockCancel()
				err = dl.TryLock(lockCtx)
				require.NoError(s.T(), err)
			},
			wantErr: dlock.ErrRefreshNotHeld,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)

			dl := tc.dlFunc(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := dl.Refresh(ctx)
			cancel()

			assert.True(t, errors.Is(err, tc.wantErr))
		})
	}
}

func (s *DlockTestSuite) TestAutoExpire() {
	t := s.T()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dl, err := s.dc.NewDlock(ctx, fmt.Sprintf("%s:autoexpire", keyPrefix), time.Second)
	require.NoError(t, err)

	lockCtx, lockCancel := context.WithTimeout(context.Background(), time.Second)
	defer lockCancel()
	err = dl.TryLock(lockCtx)
	require.NoError(t, err)

	time.Sleep(time.Second)

	newCtx, newCancel := context.WithTimeout(context.Background(), time.Second)
	defer newCancel()
	anotherDl, err := s.dc.NewDlock(newCtx, fmt.Sprintf("%s:autoexpire", keyPrefix), time.Second)
	require.NoError(t, err)
	anotherCtx, anotherCancel := context.WithTimeout(context.Background(), time.Second)
	defer anotherCancel()
	assert.NoError(t, anotherDl.TryLock(anotherCtx))

	unlockCtx, unlockCancel := context.WithTimeout(context.Background(), time.Second)
	defer unlockCancel()
	assert.True(t, errors.Is(dl.Unlock(unlockCtx), dlock.ErrReleaseNotHeld))

	refreshCtx, refreshCancel := context.WithTimeout(context.Background(), time.Second)
	defer refreshCancel()
	assert.True(t, errors.Is(dl.Refresh(refreshCtx), dlock.ErrRefreshNotHeld))
}
