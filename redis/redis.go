package redis

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/sahalazain/godilock"
	"github.com/sahalazain/godilock/redlock"
)

const schema = "redis"

// LockManager Redis lock manager
type LockManager struct {
	manager *redlock.RedLock
	prefix  string
}

func init() {
	godilock.Register(schema, New)
}

// New create redis locker instance
func New(u *url.URL) (godilock.DLocker, error) {
	host := u.Host

	hs := strings.Split(host, ",")
	for i := range hs {
		hs[i] = "tcp://" + hs[i]
	}

	lockMgr, err := redlock.NewRedLock(hs)

	if err != nil {
		return nil, err
	}

	return &LockManager{
		manager: lockMgr,
		prefix:  u.Path,
	}, nil
}

// TryLock try to lock, and return immediately if resource already locked
func (l *LockManager) TryLock(ctx context.Context, id string, ttl int) error {
	_, err := l.manager.TryLock(ctx, l.prefix+id, time.Duration(ttl)*time.Second)
	return err
}

// Lock try to lock and wait until resource is available to lock
func (l *LockManager) Lock(ctx context.Context, id string, ttl int) error {
	_, err := l.manager.Lock(ctx, l.prefix+id, time.Duration(ttl)*time.Second)
	return err
}

// Unlock unlock resource
func (l *LockManager) Unlock(ctx context.Context, id string) error {
	return l.manager.UnLock(ctx, l.prefix+id)
}

// Close close the lock
func (l *LockManager) Close() error {
	return l.manager.Close()
}
