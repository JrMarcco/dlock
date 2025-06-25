package dlock

import "errors"

var (
	ErrLockIsHeld     = errors.New("[dlock] failed to acquire lock, the lock is held by another process")
	ErrReleaseNotHeld = errors.New("[dlock] can not release the lock that not held")
	ErrRefreshNotHeld = errors.New("[dlock] can not refresh the lock that not held")
)
