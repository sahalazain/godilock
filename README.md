# godilock
Go Distributed Lock

## Usage

```
import (
	"github.com/sahalazain/godilock"
	_ "github.com/sahalazain/godilock/etcd"
	_ "github.com/sahalazain/godilock/local"
	_ "github.com/sahalazain/godilock/redis"
	_ "github.com/sahalazain/godilock/zk"
)

url := "local://"
// url := "redis://localhost:6379/test"
// url := "redis://127.0.0.1:6379,127.0.0.2:6379,127.0.0.3:6379/test"
// url := "etcd://127.0.0.1:2379,127.0.0.2:2379,127.0.0.3:2379/test"
// url := "zk://127.0.0.1:2181,127.0.0.2:2181,127.0.0.3:2181/test"

lock, _ := godilock.New(url)

// Lock and wait
if err := lock.Lock(ctx, id, 20); err != nil {
    return err
}

// Try lock and return immediately
if err := lock.TryLock(ctx, id, 20); err != nil {
    continue
}

// Unlock
if err := lock.Unlock(ctx, id); err != nil {
    return err
}

```