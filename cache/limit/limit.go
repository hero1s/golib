package limit

import (
	"fmt"
	"github.com/hero1s/golib/cache"
	"github.com/hero1s/golib/log"
	"time"
)

type AccessLimitConf struct {
	Frequency  int64 `json:"frequency"`
	ExpireTime int64 `json:"expire_time"`
}

var Al AccessLimitConf //访问限速配置

/*
 *desc:用于访问次数限制
 *@key:需要以什么来做标识做访问次数限制的标志
 *@frequency: 次数
 *@expireTime:多少秒超时
 */
func AccessLimit(key string, frequency int64, expireTime int64, isGobal bool) bool {
	if !isGobal {
		ok := cache.MemCache.IsExist(key)
		if !ok { //doesn't exist, set a key and expire time
			cache.MemCache.Put(key, 1, time.Duration(expireTime)*time.Second)
			return true
		}
		f := cache.MemCache.GetInt64(key)
		if f < frequency { // 内存模式不会修改过期时间
			cache.MemCache.Incr(key)
			return true
		}
		log.Infof(fmt.Sprintf("触及访问限速,key:%v,frequency:%v,expireTime:%v", key, frequency, expireTime))
		return false
	} else {
		ok := cache.RedisCache.IsExist(key)
		if !ok { //doesn't exist, set a key and expire time
			cache.RedisCache.Put(key, 1, time.Duration(expireTime)*time.Second)
			return true
		}
		f := cache.RedisCache.GetInt64(key)
		if f < frequency { // redis模式或修改过期时间
			cache.RedisCache.Incr(key)
			ret := cache.Redis.TTL(key)
			if ret.Val() > 0 {
				cache.RedisCache.Expire(key, time.Duration(expireTime)*time.Second-ret.Val())
			} else {
				cache.RedisCache.Expire(key, time.Duration(expireTime)*time.Second)
			}
			return true
		}
		log.Infof(fmt.Sprintf("触及全局访问限速,key:%v,frequency:%v,expireTime:%v", key, frequency, expireTime))
		return false
	}
}

//一天n次限制
func LimitDay(key string, id uint64, n int64, logLimit bool) bool {
	key = GetDateKey("d", key, id)
	re := cache.Redis.Incr(key)
	cache.Redis.Expire(key, time.Hour*24)
	if re.Err() != nil {
		log.Errorf("%v 一天n次限制,redis出错:%v", key, re.Err().Error())
		return false
	}
	log.Debugf("%v 一天N次限制:re:%v,n:%v", key, re.Val(), n)
	if re.Val() > n {
		if logLimit {
			log.Errorf("超过最大次数限制:%v --- %v", key, re.Val())
		}
		return false
	}
	return true
}

//清除每天N次限制
func LimitDayClear(key string, id uint64) {
	key = GetDateKey("d", key, id)
	cache.Redis.Del(key)
}

//获得一个按自然日期的key
//@Param time_type "d"=日,"w"=周,"m"=月
func GetDateKey(time_type string, key string, id interface{}) string {
	var str string
	if time_type == "d" {
		str = time.Now().Format("20060102")
	} else if time_type == "w" {
		y, w := time.Now().ISOWeek()
		str = fmt.Sprintf("%vW%v", y, w)
	} else if time_type == "m" {
		str = time.Now().Format("200601")
	}
	return fmt.Sprintf("%v:%v:%v", key, str, id)
}

//限制错误次数冻结
func CheckErrorLock(key string, isRight bool, limitCount, lockHour int64) bool {
	count := cache.RedisCache.GetInt64(key)
	if count > limitCount {
		log.Errorf("%v 错误次数太多:%v", key, count)
		return false
	}
	if isRight { //密码正确，清除次数
		cache.RedisCache.Delete(key)
		return true
	}
	cache.RedisCache.Incr(key)
	cache.RedisCache.Expire(key, time.Hour*time.Duration(lockHour))

	log.Infof("验证错误:%v 计数:%v", key, count+1)
	return true
}

// 获取一个key
func GetFormatKey(key string, value interface{}) string {
	return fmt.Sprintf("%v:%v", key, value)
}
