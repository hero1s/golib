package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hero1s/golib/db/orm"
	"github.com/hero1s/golib/log"
)

// 读取mysql的配置, 初始化mysql
type DbConf struct {
	AliasName string `json:"alias_name"`
	Host      string `json:"host"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DbName    string `json:"db_name"`
	DebugLog  bool   `json:"debug_log"`
	DueTime   int64  `json:"due_time"`
}

func InitDBByDsn(aliasName, dsn string, debugLog bool, dueTimeMs int64, maxIdleConn int, maxOpenConn int) error {
	orm.Debug = debugLog
	orm.LogFunc = func(queies string, err error, elsp float64) {
		if err != nil {
			log.Errorf("%v", queies)
		} else {
			if dueTimeMs <= int64(elsp) {
				log.Infof("%v", queies)
			}
		}
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	return orm.RegisterDataBase(aliasName, "mysql", dsn, maxIdleConn, maxOpenConn)
}

// init mysql params(30, 500,int64(10*time.Minute))
func InitDB(aliasName, user, password, host, dbName string, debugLog bool, dueTimeMs int64, maxIdleConn int, maxOpenConn int) error {
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local", user, password, host, dbName)
	return InitDBByDsn(aliasName, source, debugLog, dueTimeMs, maxIdleConn, maxOpenConn)
}

func InitDBConf(conf DbConf, maxIdleConn int, maxOpenConn int) error {
	err := InitDB(conf.AliasName, conf.User, conf.Password, conf.Host, conf.DbName, conf.DebugLog, conf.DueTime, maxIdleConn, maxOpenConn)
	if err != nil {
		return err
	}
	return err
}
