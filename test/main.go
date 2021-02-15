package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sahalazain/godilock"
	_ "github.com/sahalazain/godilock/etcd"
	_ "github.com/sahalazain/godilock/local"
	_ "github.com/sahalazain/godilock/redis"
	_ "github.com/sahalazain/godilock/zk"
)

var jobs map[int]bool
var wait bool

func main() {
	var mode = flag.String("m", "shared", "lock mode, shared or no")
	var url = flag.String("u", "local://", "lock url")
	var thread = flag.Int("t", 2, "number of thread")
	var iter = flag.Int("i", 10, "Number of iteration")
	flag.Parse()

	wait = true
	if *mode != "shared" {
		wait = false
	}

	var wg sync.WaitGroup
	var lock godilock.DLocker

	if *url == "local://" {
		lock, _ = godilock.New(*url)
	}

	jobs = make(map[int]bool)

	for i := 1; i <= *thread; i++ {
		wg.Add(1)
		go testLock(lock, &wg, "Thread "+fmt.Sprintf("%v", i), *url, *iter)
	}

	wg.Wait()

}

func testLock(lock godilock.DLocker, wg *sync.WaitGroup, name, surl string, iter int) {
	defer wg.Done()
	if lock == nil {
		var err error
		lock, err = godilock.New(surl)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()

	lockAcquired := 0

	for i := 0; i < iter; i++ {
		id := fmt.Sprintf("%v", i)
		if jobs[i] && !wait {
			continue
		}
		fmt.Println("["+name+"] Try to get lock for ", id)
		t := time.Now()
		if wait {
			if err := lock.Lock(ctx, id, 20); err != nil {
				fmt.Println("["+name+"] Error getting lock ", err)
				continue
			}
		} else {
			if err := lock.TryLock(ctx, fmt.Sprintf("%v", i), 10); err != nil {
				fmt.Println("["+name+"] Error getting lock ", err)
				fmt.Println("[" + name + "] Continue to next job ")
				continue
			}
		}
		d := time.Now().Sub(t)
		fmt.Println("["+name+"] Acquired lock for ", id, " in ", d.Milliseconds(), " ms")
		lockAcquired++
		jobs[i] = true
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		t = time.Now()
		if err := lock.Unlock(ctx, id); err != nil {
			fmt.Println("["+name+"] Error unlocking", err)
		}
		d = time.Now().Sub(t)
		fmt.Println("["+name+"] Unlock for ", id, " in ", d.Milliseconds(), " ms")
	}

	fmt.Println("["+name+"] Total lock acquired ", lockAcquired)
}
