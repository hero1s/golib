package db

import (
	"fmt"
	"github.com/hero1s/golib/db/dbsql"
	"github.com/hero1s/golib/db/orm"
)

type Table struct {
	TableName string
	DbName    string
}

func (t *Table) NewOrm(multiOrm ...orm.Ormer) orm.Ormer {
	o := dbsql.NewOrm(multiOrm, t.DbName)
	return o
}

func (t *Table) NewTableRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewTableRecord(t.DbName, t.TableName, data, multiOrm...)
}

func (t *Table) NewOrUpdateRecord(data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewOrUpdateRecord(t.DbName, t.TableName, data, multiOrm...)
}

func (t *Table) NewOrUpdateByAddRecord(addData map[string]interface{}, upData map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.NewOrUpdateByAddRecord(t.DbName, t.TableName, addData, upData, multiOrm...)
}

func (t *Table) UpdateTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.UpdateTableRecord(t.DbName, t.TableName, data, conditionSql, multiOrm...)
}

func (t *Table) UpdateByAddTableRecord(data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.UpdateByAddTableRecord(t.DbName, t.TableName, data, conditionSql, multiOrm...)
}

func (t *Table) DeleteTableRecord(conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.DeleteTableRecord(t.DbName, t.TableName, conditionSql, multiOrm...)
}

// 获取单行记录
func (t *Table) SingleRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	return dbsql.SingleRecordByAny(t.DbName, t.TableName, conditionSql, record, multiOrm...)
}

// 获取多行记录
func (t *Table) MultiRecordByAny(conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	return dbsql.MultiRecordByAny(t.DbName, t.TableName, conditionSql, record, multiOrm...)
}

// 带分页的多行记录并返回总行数
func (t *Table) MultiRecordByAnyOrderLimit(conditionSql string, orderbyCondition string, pageIndex, pageSize int64, record interface{}, multiOrm ...orm.Ormer) (int64, error) {
	return dbsql.MultiRecordByAnyOrderLimit(t.DbName, t.TableName, conditionSql, orderbyCondition, pageIndex, pageSize, record)
}

func (t *Table) ALLRecord(record interface{}, multiOrm ...orm.Ormer) error {
	return t.MultiRecordByAny("", record, multiOrm...)
}

// 统计条件数量
func (t *Table) CountRecord(conditionSql string) int64 {
	sql := fmt.Sprintf(`SELECT  COUNT(1) FROM %v WHERE %v`, t.TableName, conditionSql)
	var count int64
	err := t.NewOrm().Raw(sql).QueryRow(&count)
	if err != nil {
		return 0
	}
	return count
}
