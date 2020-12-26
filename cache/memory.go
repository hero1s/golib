
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// DefaultEvery means the clock time of recycling the expired cache items in memory.
	DefaultEvery = 60 // 1 minute
)

// MemoryItem store memory cache item.
type MemoryItem struct {
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Now().Sub(mi.createdTime) > mi.lifespan
}

func (mi *MemoryItem) setExpire(lifespan time.Duration) {
	mi.createdTime = time.Now()
	mi.lifespan = lifespan
}

// MemoryCache is Memory cache adapter.
// it contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) Get(name string) interface{} {
	bc.RLock()
	defer bc.RUnlock()
	if itm, ok := bc.items[name]; ok {
		if itm.isExpire() {
			return nil
		}
		return itm.val
	}
	return nil
}

type Error string

func (err Error) Error() string { return string(err) }

var ErrNil = errors.New("memory: nil returned")

// GetString is empty function
func (bc *MemoryCache) GetString(name string) string {
	val := bc.Get(name)
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func (bc *MemoryCache) GetInt(name string) int {
	val := bc.Get(name)
	if val == nil {
		return 0
	}
	value := fmt.Sprintf("%v", val)
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int(n)
}

func (bc *MemoryCache) GetInt64(name string) int64 {
	val := bc.Get(name)
	if val == nil {
		return 0
	}
	value := fmt.Sprintf("%v", val)
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func (bc *MemoryCache) GetFloat64(name string) float64 {
	val := bc.Get(name)
	if val == nil {
		return 0
	}
	value := fmt.Sprintf("%v", val)
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return n
}

func (bc *MemoryCache) GetBool(name string) bool {
	val := bc.Get(name)
	if val == nil {
		return false
	}
	value := fmt.Sprintf("%v", val)
	ok, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return ok
}

// GetMulti gets caches from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) GetMulti(names []string) []interface{} {
	var rc []interface{}
	for _, name := range names {
		rc = append(rc, bc.Get(name))
	}
	return rc
}

// Put cache to memory.
// if lifespan is 0, it will be forever till restart.
func (bc *MemoryCache) Put(name string, value interface{}, lifespan time.Duration) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[name] = &MemoryItem{
		val:         value,
		createdTime: time.Now(),
		lifespan:    lifespan,
	}
	return nil
}

// Put cache to memory.
// if lifespan is 0, it will be forever till restart.
func (bc *MemoryCache) Expire(name string, lifespan time.Duration) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[name].setExpire(lifespan)
	return nil
}

// Delete cache in memory.
func (bc *MemoryCache) Delete(name string) error {
	bc.Lock()
	defer bc.Unlock()
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// Incr increase cache counter in memory.
// it supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) Incr(key string) error {
	bc.RLock()
	defer bc.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) + 1
	case int32:
		itm.val = itm.val.(int32) + 1
	case int64:
		itm.val = itm.val.(int64) + 1
	case uint:
		itm.val = itm.val.(uint) + 1
	case uint32:
		itm.val = itm.val.(uint32) + 1
	case uint64:
		itm.val = itm.val.(uint64) + 1
	default:
		return errors.New("item val is not (u)int (u)int32 (u)int64")
	}
	return nil
}

// Decr decrease counter in memory.
func (bc *MemoryCache) Decr(key string) error {
	bc.RLock()
	defer bc.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) - 1
	case int64:
		itm.val = itm.val.(int64) - 1
	case int32:
		itm.val = itm.val.(int32) - 1
	case uint:
		if itm.val.(uint) > 0 {
			itm.val = itm.val.(uint) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint32:
		if itm.val.(uint32) > 0 {
			itm.val = itm.val.(uint32) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint64:
		if itm.val.(uint64) > 0 {
			itm.val = itm.val.(uint64) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
}

// IsExist check cache exist in memory.
func (bc *MemoryCache) IsExist(name string) bool {
	bc.RLock()
	defer bc.RUnlock()
	if v, ok := bc.items[name]; ok {
		return !v.isExpire()
	}
	return false
}

// ClearAll will delete all cache in memory.
func (bc *MemoryCache) ClearAll() error {
	bc.Lock()
	defer bc.Unlock()
	bc.items = make(map[string]*MemoryItem)
	return nil
}

// StartAndGC start memory cache. it will check expiration in every clock time.
func (bc *MemoryCache) StartAndGC(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur := time.Duration(cf["interval"]) * time.Second
	bc.Every = cf["interval"]
	bc.dur = dur
	go bc.vaccuum()
	return nil
}

// check expiration.
func (bc *MemoryCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	for {
		<-time.After(bc.dur)
		if bc.items == nil {
			return
		}
		if keys := bc.expiredKeys(); len(keys) != 0 {
			bc.clearItems(keys)
		}
	}
}

// expiredKeys returns key list which are expired.
func (bc *MemoryCache) expiredKeys() (keys []string) {
	bc.RLock()
	defer bc.RUnlock()
	for key, itm := range bc.items {
		if itm.isExpire() {
			keys = append(keys, key)
		}
	}
	return
}

// clearItems removes all the items which key in keys.
func (bc *MemoryCache) clearItems(keys []string) {
	bc.Lock()
	defer bc.Unlock()
	for _, key := range keys {
		delete(bc.items, key)
	}
}

func init() {
	Register("memory", NewMemoryCache)
}
