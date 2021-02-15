package local

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/sahalazain/godilock"
)

const schema = "local"
const interval = 100

// LockManager local lock manager
type LockManager struct {
	mux    *sync.Mutex
	locked map[string]bool
}

func init() {
	godilock.Register(schema, New)
}

// New create redis locker instance
func New(u *url.URL) (godilock.DLocker, error) {
	return &LockManager{
		mux:    &sync.Mutex{},
		locked: make(map[string]bool),
	}, nil
}

func (l *LockManager) lock(id string, ttl int) error {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.locked[id] {
		return errors.New("resource locked")
	}

	l.locked[id] = true

	go func(l *LockManager) {
		time.Sleep(time.Duration(ttl) * time.Second)
		l.Unlock(nil, id)
	}(l)

	return nil
}

// TryLock try to lock, and return immediately if resource already locked
func (l *LockManager) TryLock(ctx context.Context, id string, ttl int) error {
	return l.lock(id, ttl)
}

// Lock try to lock and wait until resource is available to lock
func (l *LockManager) Lock(ctx context.Context, id string, ttl int) error {
	err := l.lock(id, ttl)
	if err == nil {
		return nil
	}

	count := 0
	max := ttl * 1000 / interval
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		err := l.lock(id, ttl)
		if err == nil {
			return nil
		}
		count++
		if count > max {
			return err
		}
	}
}

// Unlock unlock resource
func (l *LockManager) Unlock(ctx context.Context, id string) error {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.locked[id] = false
	return nil
}

// Close close the lock
func (l LockManager) Close() error {
	return nil
}
