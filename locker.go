package godilock

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

//DLocker distributed locker interface
type DLocker interface {
	TryLock(ctx context.Context, id string, ttl int) error
	Lock(ctx context.Context, id string, ttl int) error
	Unlock(ctx context.Context, id string) error
	Close() error
}

//InitFunc cache init function
type InitFunc func(url *url.URL) (DLocker, error)

var lockerImpl = make(map[string]InitFunc)

//Register register cache implementation
func Register(schema string, f InitFunc) {
	lockerImpl[schema] = f
}

//New create new cache
func New(urlStr string) (DLocker, error) {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	f, ok := lockerImpl[u.Scheme]
	if !ok {
		return nil, errors.New("Unsupported scheme")
	}

	return f(u)
}
