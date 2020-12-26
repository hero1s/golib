package uuid

import (
	"git.moumentei.com/plat_go/golib/log"
	"git.moumentei.com/plat_go/golib/utils/uuid/internal"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"strconv"
)

var (
	Uid *internal.UUID
)

func InitUUID(redisHost, password string) error {
	newClient := func() (redis.Cmdable, bool, error) {
		return redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: password,
		}), true, nil
	}
	Uid = internal.NewUUID("uid")
	err := Uid.LoadH52FromRedis(newClient, "UUID:UID:24")
	if err != nil {
		log.Error("初始化UUID错误:%v", err.Error())
		return err
	}

	return nil
}

func GenUid() uint64 {
	return Uid.Next()
}

func GenStringUUID() string {
	return strconv.FormatUint(GenUid(), 10)
}

func NewUuid() string {
	return uuid.New().String()
}
