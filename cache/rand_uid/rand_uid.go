package rand_uid

import (
	"fmt"
	"github.com/hero1s/golib/cache"
	"github.com/hero1s/golib/db/orm"
	"github.com/hero1s/golib/helpers/math"
	"github.com/hero1s/golib/log"
	"gopkg.in/fatih/set.v0"
)

//范围随机ID生成(例如生成不可预测的UID)

var randUidKey = "rand_uid_key"

//初始化
func InitRandUid(key string) {
	randUidKey = key
}

//获取一个uid
func PopUid() (int64, error) {
	res := cache.Redis.SPop(randUidKey)
	uid, err := res.Int64()
	if err != nil {
		log.Errorf("获取uid错误:%v", err)
	}
	reslen := cache.Redis.SCard(randUidKey)
	log.Infof("获取用户ID:%v,剩余ID数量:%v", uid, reslen.Val())
	return uid, err
}

//重新生成uid
func ResetNewUid(startId, endId, num int64) {
	a := set.New(set.ThreadSafe)
	for i := int64(0); i < num; i = i + 1 {
		id := math.Random(startId, endId)
		a.Add(id)
	}
	a.Each(func(id interface{}) bool {
		cache.Redis.SAdd(randUidKey, id)
		return true
	})
	res := cache.Redis.SCard(randUidKey)
	log.Infof("重新生成用户ID:%v-%v:生成数量:%v,剩余数量:%v", startId, endId, num, res.Val())
}

//通过数据库表字段重新生成
func ResetNewUidByTable(tableName, field string, step, addNum, getNum int64, o orm.Ormer) error {
	var curMaxId int64
	sql := fmt.Sprintf("SELECT MAX(%v) FROM %v", field, tableName)
	err := o.Raw(sql).QueryRow(&curMaxId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	ResetNewUid(curMaxId+step, curMaxId+addNum, getNum)
	return nil
}
