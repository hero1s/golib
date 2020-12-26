package dbsql

import (
	"fmt"
	"time"
)

//分表操作(按月分)

// 使用者可能要在前面封装一层
func CreateNextMonthTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		CreateThisMonthTable(dbName, oldTable, newTable)()
		t := GetMonthTableByTimestamp(newTable, time.Now().AddDate(0, 1, 0).Unix())
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

func CreateThisMonthTable(dbName, oldTable, newTable string) func() error {
	return func() error {
		t := GetMonthTableByTimestamp(newTable, time.Now().Unix())
		sql := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %v LIKE %v
		`, t, oldTable)

		o := NewOrmWithDB(dbName)
		_, err := o.Raw(sql).Exec()
		return err
	}
}

//获取今天所在的月份表名
func GetTodayMonthTable(tableName string) string {
	return GetMonthTableByTimestamp(tableName, time.Now().Unix())
}

//获取N个月前后的表名
func GetDiffMonthTable(tableName string, diffMonth int) string {
	return GetMonthTableByTimestamp(tableName, time.Now().AddDate(0, diffMonth, 0).Unix())
}

//根据时间戳获取格式华的日期
func GetMonthTableByTimestamp(tableName string, timestamp int64) string {
	return tableName + time.Unix(timestamp, 0).Format("200601")
}
