package stimer

import (
	"github.com/hero1s/golib/cache"
	"github.com/hero1s/golib/log"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

//基于redis的分布式定时器

const (
	TimerBaseName = "STimers:"
)

type CallFunc func(tid uint64)

var timers = map[string]CallFunc{}

func InitTimerByName(name string, callFunc CallFunc) {
	name = TimerBaseName + name
	timers[name] = callFunc
}

func Tick() {
	for k, v := range timers {
		ids := popTimer(k)
		for _, id := range ids {
			v(id)
			log.Infof("定时器:%v 超时ids:%v", k, id)
		}
	}
}

//添加
func AddTimer(name string, tid uint64, sec int64) {
	name = TimerBaseName + name
	if _, ok := timers[name]; !ok {
		log.Errorf("%v 定时器未初始化", name)
		return
	}
	score := time.Now().Unix() + sec
	re := cache.Redis.ZAdd(name, &redis.Z{Score: float64(score), Member: tid})
	if re.Err() != nil {
		log.Errorf("设置定时有序队列出错:%v", re.Err().Error())
		return
	}
	log.Infof("%v 添加定时器:%v,%v", name, tid, sec)
}

//返回超时定时器id
func popTimer(name string) []uint64 {
	ids := make([]uint64, 0)
	//取出定时器排行傍数据
	now := strconv.FormatInt(time.Now().Unix(), 10)
	re := cache.Redis.ZRangeByScore(name, &redis.ZRangeBy{Min: "0", Max: now, Offset: 0, Count: 100})
	if re.Err() != nil || len(re.Val()) == 0 {
		return ids
	}
	for _, v := range re.Val() {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			ids = append(ids, id)
		}
	}
	for _, v := range ids {
		cache.Redis.ZRem(name, v)
	}
	return ids
}
func DelTimer(name string, tid uint64) {
	name = TimerBaseName + name
	log.Infof("%v 删除定时器:%v", name, tid)
	cache.Redis.ZRem(name, tid)
}
