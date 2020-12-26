package dbsql

import (
	"errors"
	"github.com/hero1s/golib/db/orm"
	"github.com/hero1s/golib/i18n"
)

func NewOrmWithDB(db string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(db)
	return o
}

func NewOrm(multiOrm []orm.Ormer, db ...string) (o orm.Ormer) {
	if len(multiOrm) == 0 {
		o = orm.NewOrm()
		if len(db) == 1 {
			o.Using(db[0])
		}
	} else if len(multiOrm) == 1 {
		o = multiOrm[0]
	} else {
		panic("只能传一个Ormer")
	}
	return
}

//构造查询
func NewQueryBuilder() (qb orm.QueryBuilder, err error) {
	qb, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, i18n.WrapDatabaseError(err)
	}

	return qb, nil
}

func NewTableRecord(dbName, tableName string, data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := InsertSql(tableName, data)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.LastInsertId()
	return n, i18n.WrapDatabaseError(err)
}

// 一次性插入多条值
func NewMultiTableRecord(dbName, tableName string, data map[string][]interface{}, multiOrm ...orm.Ormer) error {
	// filterData(data, t.TableField)
	values, sql := MultiInsertSql(tableName, data)
	if len(values) == 0 {
		return i18n.WrapDatabaseError(errors.New("insert multi的value数据不齐"))
	}
	o := NewOrm(multiOrm, dbName)
	rp, err := o.Raw(sql).Prepare()
	defer rp.Close()
	if err != nil {
		return i18n.WrapDatabaseError(err)
	}
	for _, value := range values {
		_, err := rp.Exec(value...)
		if err != nil {
			return i18n.WrapDatabaseError(err)
		}
	}
	return rp.Close()
}

func NewOrUpdateRecord(dbName, tableName string, data map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := InsertOrUpdateSql(tableName, data)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.LastInsertId()
	return n, i18n.WrapDatabaseError(err)
}

func NewOrUpdateByAddRecord(dbName, tableName string, addData map[string]interface{}, upData map[string]interface{}, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := InsertOrUpdateByAddSql(tableName, addData, upData)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.LastInsertId()
	return n, i18n.WrapDatabaseError(err)
}

func UpdateTableRecord(dbName, tableName string, data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := UpdateSql(tableName, data, conditionSql)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

func UpdateByAddTableRecord(dbName, tableName string, data map[string]interface{}, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	// filterData(data, t.TableField)
	o := NewOrm(multiOrm, dbName)
	values, sql := UpdateByAddSql(tableName, data, conditionSql)
	result, err := o.Raw(sql, values).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

func DeleteTableRecord(dbName, tableName string, conditionSql string, multiOrm ...orm.Ormer) (int64, error) {
	o := NewOrm(multiOrm, dbName)
	sql := DeleteSql(tableName, conditionSql)
	result, err := o.Raw(sql).Exec()
	if err != nil {
		return 0, i18n.WrapDatabaseError(err)
	}
	n, err := result.RowsAffected()
	return n, i18n.WrapDatabaseError(err)
}

// 获取单行记录
func SingleRecordByAny(dbName, tableName string, conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	o := NewOrm(multiOrm, dbName)
	sql := QuerySingleSql(tableName, conditionSql)
	err := o.Raw(sql).QueryRow(record)
	if err == orm.ErrNoRows {
		return i18n.RecordNotFound
	}
	return i18n.WrapDatabaseError(err)
}

//获取多行记录
func MultiRecordByAny(dbName, tableName string, conditionSql string, record interface{}, multiOrm ...orm.Ormer) error {
	o := NewOrm(multiOrm, dbName)
	sql := QueryMultiSql(tableName, conditionSql)
	_, err := o.Raw(sql).QueryRows(record)
	return i18n.WrapDatabaseError(err)
}

//多行记录带分页,并返回总行数
func MultiRecordByAnyOrderLimit(dbName, tableName string, conditionSql string, orderbyCondition string, pageIndex, pageSize int64, record interface{}, multiOrm ...orm.Ormer) (int64, error) {
	o := NewOrm(multiOrm, dbName)
	sql := QueryMultiSqlOderByLimit(tableName, conditionSql, orderbyCondition, pageIndex, pageSize)
	return MultiRecordAndTotal(sql, record, o)
}

//多行记录带分页,并返回总行数
func MultiRecordAndTotal(querySql string, record interface{}, o orm.Ormer) (int64, error) {
	totalSql := QueryLimitTotalSql(querySql)
	_, err := o.Raw(querySql).QueryRows(record)
	var total int64
	if err != nil {
		return total, i18n.WrapDatabaseError(err)
	}
	o.Raw(totalSql).QueryRow(&total)
	return total, nil
}

func ALLRecord(dbName, tableName string, record interface{}, multiOrm ...orm.Ormer) error {
	return MultiRecordByAny(dbName, tableName, "", record, multiOrm...)
}

//返回true就表示不存在
func CheckNoExist(err error) bool {
	if err != nil && err == orm.ErrNoRows {
		return true
	}
	return false
}

func WrapDatabaseError(err error) error {
	if err == orm.ErrNoRows {
		return i18n.RecordNotFound
	}
	return i18n.WrapDatabaseError(err)
}
