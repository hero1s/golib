
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "car"
)

// Cache is Redis cache adapter.
type RedisItem struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() Cache {
	//return &Cache{key: DefaultKey}
	return &RedisItem{}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *RedisItem) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *RedisItem) associate(originKey interface{}) string {
	//return fmt.Sprintf("%s-%s", rc.key, originKey)
	return fmt.Sprintf("%s", originKey)
}

// Get cache from redis.
func (rc *RedisItem) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetString is empty function
func (rc *RedisItem) GetString(key string) string {
	if v, err := redis.String(rc.do("GET", key)); err == nil {
		return v
	}
	return ""
}

func (rc *RedisItem) GetInt(key string) int {
	if v, err := redis.Int(rc.do("GET", key)); err == nil {
		return v
	}
	return 0
}

func (rc *RedisItem) GetInt64(key string) int64 {
	if v, err := redis.Int64(rc.do("GET", key)); err == nil {
		return v
	}
	return 0
}

func (rc *RedisItem) GetFloat64(key string) float64 {
	if v, err := redis.Float64(rc.do("GET", key)); err == nil {
		return v
	}
	return 0
}

func (rc *RedisItem) GetBool(key string) bool {
	if v, err := redis.Bool(rc.do("GET", key)); err == nil {
		return v
	}
	return false
}

// GetMulti get cache from redis.
func (rc *RedisItem) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *RedisItem) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Put put cache to redis.
func (rc *RedisItem) Expire(key string, timeout time.Duration) error {
	_, err := rc.do("EXPIRE", key, int64(timeout/time.Second))
	return err
}

// Delete delete cache in redis.
func (rc *RedisItem) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *RedisItem) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *RedisItem) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *RedisItem) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *RedisItem) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	//cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+"-*"))
	cachedKeys, err := redis.Strings(c.Do("KEYS", "*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisItem) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *RedisItem) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	Register("redis", NewRedisCache)
}
