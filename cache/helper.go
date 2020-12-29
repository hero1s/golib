package cache

import (
	"github.com/hero1s/golib/constant"
	"github.com/hero1s/golib/log"
	"time"
)

// 设置db缓存
func SetDBCache(key string, value interface{}, timeout time.Duration) error {
	if timeout == 0 {
		timeout = constant.SECOND_ONE_DAY * time.Second
	}
	return SetCache(RedisCache, key, value, timeout)
}
func DeleteDBCache(key string) error {
	return DelCache(RedisCache, key)
}
func GetDBCache(key string, to interface{}) error {
	return GetCache(RedisCache, key, to)
}

// 从缓存或DB中取数据
// value 传引用
func GetDataFromCacheOrDB(key string, value interface{}, forceSql bool, timeout time.Duration, dbFunc func(value interface{}) error) error {
	if forceSql {
		err := dbFunc(value)
		if err == nil {
			log.Debugf("强制设置缓存:%v", key)
			SetDBCache(key, value, timeout)
			return nil
		}
		log.Debugf("%v获取数据失败:%v", key, err)
		return err
	} else {
		err := GetDBCache(key, value)
		if err == nil {
			log.Debugf("缓存获取数据成功:%v", key)
			return nil
		} else {
			err := dbFunc(value)
			if err == nil {
				log.Debugf("数据库获取数据设置缓存:%v", key)
				SetDBCache(key, value, timeout)
				return nil
			}
			log.Errorf("%v获取数据失败:%v", key, err)
			return err
		}
	}
}
