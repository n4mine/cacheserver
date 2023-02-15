package cache

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/n4mine/cacheserver/chunks"
	"github.com/n4mine/cacheserver/config"
	"github.com/n4mine/cacheserver/models"
)

const SHARD_COUNT = 32

var CacheObj caches

var (
	TotalCount int64
	cleaning   bool
)

type (
	caches []*cache
)

type cache struct {
	items map[string]*chunks.CS // [counter]ts,value
	sync.RWMutex
}

func InitCaches() {
	CacheObj = NewCaches()
}

func NewCaches() caches {
	c := make(caches, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		c[i] = &cache{items: make(map[string]*chunks.CS)}
	}
	return c
}

func GC(cfg config.GcConfig) {
	t := time.NewTicker(time.Minute * time.Duration(cfg.GcIntervalInMinutes))
	cleaning = false

	for {
		select {
		case <-t.C:
			if !cleaning {
				go CacheObj.GC(cfg.ExpiresInMinutes)
			} else {
				log.Println("GC() is working, may be it's too slow")
			}
		}
	}
}

func (c *caches) Push(seriesID string, ts int64, value float64) error {
	shard := c.getShard(seriesID)
	existC, exist := CacheObj.exist(seriesID)

	if exist {
		shard.Lock()
		err := existC.Push(ts, value)
		shard.Unlock()
		return err
	}
	newC := CacheObj.create(seriesID)
	shard.Lock()
	err := newC.Push(ts, value)
	shard.Unlock()

	return err
}

func (c *caches) Get(seriesID string, from, to int64) ([]chunks.Iter, error) {
	existC, exist := CacheObj.exist(seriesID)

	if !exist {
		return nil, models.ErrNonExistSeries
	}

	// f(series's oldestTs) ... t(series's newestTs)
	f, t := existC.GetInfo()
	// 为什么没有 to > int64(t)
	// cacheserver的数据都是经由transfer进行时间戳对齐的, 大部分场景都不会满足 to > t
	if from < int64(f) || from > int64(t) || to < int64(f) {
		return nil, models.ErrNonEnoughData
	}

	res := existC.Get(from, to)
	if res == nil {
		return nil, models.ErrInternalError
	}

	return res, nil
}

func (c *caches) create(seriesID string) *chunks.CS {
	atomic.AddInt64(&TotalCount, 1)
	shard := c.getShard(seriesID)
	shard.Lock()
	newC := chunks.NewChunks(config.C.Cache)
	shard.items[seriesID] = newC
	shard.Unlock()

	return newC
}

func (c *caches) exist(seriesID string) (*chunks.CS, bool) {
	shard := c.getShard(seriesID)
	shard.RLock()
	existC, exist := shard.items[seriesID]
	shard.RUnlock()

	return existC, exist
}

func (c caches) Count() int64 {
	return atomic.LoadInt64(&TotalCount)
}

func (c caches) remove(seriesID string) {
	atomic.AddInt64(&TotalCount, -1)
	shard := c.getShard(seriesID)
	shard.Lock()
	delete(shard.items, seriesID)
	shard.Unlock()
}

func (c caches) GC(expiresInMinutes int) {
	now := time.Now()
	done := make(chan struct{})
	var count int64
	cleaning = true
	defer func() { cleaning = false }()

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)

		for _, shard := range c {
			go func(shard *cache) {
				shard.RLock()
				for key, chunks := range shard.items {
					_, lastTs := chunks.GetInfoUnsafe()
					if int64(lastTs) < now.Unix()-60*int64(expiresInMinutes) {
						atomic.AddInt64(&count, 1)
						shard.RUnlock()
						c.remove(key)
						shard.RLock()
					}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		done <- struct{}{}
	}()

	<-done
	log.Printf("cleanup %v items, took %.2f ms\n", count, float64(time.Since(now).Nanoseconds())*1e-6)
}

func (c caches) GetInfoByKey(key string) map[string]uint32 {
	existC, exist := CacheObj.exist(key)
	if exist {
		oldest, newest := existC.GetInfo()
		return map[string]uint32{"oldest": oldest, "newest": newest, "duration": newest - oldest}
	}
	return map[string]uint32{"oldest": 0, "newest": 0, "duration": 0}
}

func (c caches) getShard(key string) *cache {
	return c[fnv32(key)%uint32(SHARD_COUNT)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
