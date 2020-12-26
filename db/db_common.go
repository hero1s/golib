package db

import (
	"fmt"
	"github.com/hero1s/golib/db/dbsql"
	"github.com/hero1s/golib/db/orm"
)

type Table interface {
	TableName() string
	DbName() string
	FullTableName() string
}

type TableOper struct {
	Tb Table
}

func (t *TableOper) NewOrm(multiOrm ...orm.Ormer) orm.Ormer {
	o := dbsql.NewOrm(multiOrm, t.Tb.DbName())
	return o
}

func (t *TableOper) NewTableRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewTableRecord(t.Tb.DbName(), t.Tb.TableName(), data, multiOrm...)
}

func (t *TableOper) NewOrUpdateRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewOrUpdateRecord(t.Tb.DbName(), t.Tb.TableName(), data, multiOrm...)
}

func (t *TableOper) NewOrUpdateByAddRecord(addData map[string]interface{}, upData map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewOrUpdateByAddRecord(t.Tb.DbName(), t.Tb.TableName(), addData, upData, multiOrm...)
}

func (t *TableOper) UpdateTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.UpdateTableRecord(t.Tb.DbName(), t.Tb.TableName(), data, conditionSql, multiOrm...)
}

func (t *TableOper) UpdateByAddTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.UpdateByAddTableRecord(t.Tb.DbName(), t.Tb.TableName(), data, conditionSql, multiOrm...)
}

func (t *TableOper) DeleteTableRecord(conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.DeleteTableRecord(t.Tb.DbName(), t.Tb.TableName(), conditionSql, multiOrm...)
}

// 获取单行记录
func (t *TableOper) SingleRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	return dbsql.SingleRecordByAny(t.Tb.DbName(), t.Tb.TableName(), conditionSql, record, multiOrm...)
}

// 获取多行记录
func (t *TableOper) MultiRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	return dbsql.MultiRecordByAny(t.Tb.DbName(), t.Tb.TableName(), conditionSql, record, multiOrm...)
}

// 带分页的多行记录并返回总行数
func (t *TableOper) MultiRecordByAnyOrderLimit(conditionSql string, orderbyCondition string, pageIndex, pageSize int64, record interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.MultiRecordByAnyOrderLimit(t.Tb.DbName(), t.Tb.TableName(), conditionSql, orderbyCondition, pageIndex, pageSize, record)
}

func (t *TableOper) ALLRecord(record interface{}, multiOrm ...orm.Ormer) error {
	return t.MultiRecordByAny("", record, multiOrm...)
}

// 统计条件数量
func (t *TableOper) CountRecord(conditionSql string) int64 {
	sql := fmt.Sprintf(`SELECT  COUNT(1) FROM %v WHERE %v`, t.Tb.TableName(), conditionSql)
	var count int64
	err := t.NewOrm().Raw(sql).QueryRow(&count)
	if err != nil {
		return 0
	}
	return count
}
