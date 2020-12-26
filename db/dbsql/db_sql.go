package dbsql

import (
	"fmt"
	"strings"
)

// 拼接update的sql语句
func UpdateSql(tableName string, data map[string]interface{}, condition string) (values []interface{}, sql string) {
	field := ""
	for k, v := range data {
		if field == "" {
			field = fmt.Sprintf("`%v`", k) + "=?"
		} else {
			field = field + "," + fmt.Sprintf("`%v`", k) + "=?"
		}
		values = append(values, v)
	}
	sql = "UPDATE " + tableName + " SET " + field + " WHERE " + condition
	return
}

func UpdateByAddSql(tableName string, data map[string]interface{}, condition string) (values []interface{}, sql string) {
	field := ""
	for k, v := range data {
		if field == "" {
			field = fmt.Sprintf("`%v`=`%v`+?", k, k)
		} else {
			field = field + "," + fmt.Sprintf("`%v`=`%v`+?", k, k)
		}
		values = append(values, v)
	}
	sql = "UPDATE " + tableName + " SET " + field + " WHERE " + condition
	return
}

func DeleteSql(tableName string, condition string) string {
	return "DELETE FROM " + tableName + " WHERE " + condition
}

func InsertSql(tableName string, data map[string]interface{}) (values []interface{}, sql string) {
	field := ""
	palceholder := ""
	for k, v := range data {
		if field != "" {
			field = field + "," + fmt.Sprintf("`%v`", k)
			palceholder = palceholder + "," + "?"
		} else {
			field = fmt.Sprintf("`%v`", k)
			palceholder = "?"
		}
		values = append(values, v)
	}
	sql = "INSERT INTO " + tableName + "(" + field + ") VALUES(" + palceholder + ")"
	return
}

func MultiInsertSql(tableName string, data map[string][]interface{}) (values [][]interface{}, sql string) {
	field := ""
	palceholder := ""
	var tmp [][]interface{}
	var valueLength int
	for k, v := range data {
		if valueLength == 0 {
			valueLength = len(v)
		} else {
			if valueLength != len(v) {
				return nil, ""
			}
		}
		if field != "" {
			field = field + "," + fmt.Sprintf("`%v`", k)
			palceholder = palceholder + "," + "?"
		} else {
			field = fmt.Sprintf("`%v`", k)
			palceholder = "?"
		}
		tmp = append(tmp, v)
	}
	for i := 0; i < valueLength; i++ {
		var value []interface{}
		for _, v := range tmp {
			value = append(value, v[i])
		}
		values = append(values, value)
	}

	sql = "INSERT INTO " + tableName + "(" + field + ") VALUES(" + palceholder + ")"
	return
}

func InsertOrUpdateByAddSql(tableName string, addData map[string]interface{}, upData map[string]interface{}) (values []interface{}, sql string) {
	values, sql = InsertSql(tableName, addData)
	sql = sql + " ON DUPLICATE KEY UPDATE "

	field := ""
	for k, v := range upData {
		if field == "" {
			field = fmt.Sprintf("`%v`=`%v`+?", k, k)
		} else {
			field = field + "," + fmt.Sprintf("`%v`=`%v`+?", k, k)
		}
		values = append(values, v)
	}
	sql = sql + field
	return
}

func InsertOrUpdateSql(tableName string, data map[string]interface{}) (values []interface{}, sql string) {
	values, sql = InsertSql(tableName, data)
	sql = sql + " ON DUPLICATE KEY UPDATE "

	field := ""
	for k, v := range data {
		if field != "" {
			field = field + ", " + fmt.Sprintf("`%v`", k) + "=?"
		} else {
			field = fmt.Sprintf("`%v`", k) + "=?"
		}
		values = append(values, v)
	}
	sql = sql + field
	return
}

func QuerySingleSql(tableName string, condition string) string {
	if condition == "" {
		return "SELECT * FROM " + tableName + " LIMIT 1"
	}
	return "SELECT * FROM " + tableName + " WHERE " + condition + " LIMIT 1"
}

func QueryMultiSql(tableName string, condition string) string {
	if condition == "" {
		return "SELECT * FROM " + tableName
	}
	return "SELECT * FROM " + tableName + " WHERE " + condition
}

func QueryMultiSqlOderByLimit(tableName string, condition string, orderbyCondition string, pageIndex, pageSize int64) string {
	pageIndex, pageSize = GetLimitAndPageSize(pageIndex, pageSize)
	if condition == "" {
		return "SELECT * FROM " + tableName + fmt.Sprintf(" order by %v limit %v,%v", orderbyCondition, pageIndex, pageSize)
	}
	return "SELECT * FROM " + tableName + " WHERE " + condition + fmt.Sprintf(" order by %v limit %v,%v", orderbyCondition, pageIndex, pageSize)
}

//添加系统函数统计影响函数
func QueryLimitTotalSql(sql string) string {
	sql = strings.ReplaceAll(sql, "FROM", "from")
	sql = "select count(*) " + sql[strings.Index(sql, "from"):]
	sql = strings.ReplaceAll(sql, "LIMIT", "limit")
	sql = sql[0:strings.Index(sql, "limit")]
	return sql
}

func FormatLimit(fileName string, pageIndex, pageSize int64) string {
	pageIndex, pageSize = GetLimitAndPageSize(pageIndex, pageSize)
	if len(fileName) > 0 {
		return fmt.Sprintf(" order by %v limit %v,%v", fileName, pageIndex, pageSize)
	} else {
		return fmt.Sprintf(" limit %v,%v", pageIndex, pageSize)
	}
}

const (
	DefaultPageSize  = 10
	MaxPageSize      = 1000
	DefaultPageIndex = 0
	MaxPageIndex     = 1000
)

func GetLimitAndPageSize(params ...int64) (int64, int64) {
	if len(params) >= 2 {
		pageIndex, pageSize := params[0], params[1]
		if pageIndex <= 0 {
			pageIndex = DefaultPageIndex
		}
		if pageIndex > MaxPageIndex {
			pageIndex = MaxPageIndex
		}
		if pageSize <= 0 {
			pageSize = DefaultPageSize
		}
		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
		return (pageIndex) * pageSize, pageSize
	}
	return 0, DefaultPageSize
}

func GetLimitAndPageSizeByMaxSize(maxIndex, maxSize int64, params ...int64) (int64, int64) {
	if len(params) >= 2 {
		pageIndex, pageSize := params[0], params[1]
		if pageIndex <= 0 {
			pageIndex = DefaultPageIndex
		}
		if pageIndex > maxIndex {
			pageIndex = maxIndex
		}
		if pageSize <= 0 {
			pageSize = DefaultPageSize
		}
		if pageSize > maxSize {
			pageSize = maxSize
		}
		return (pageIndex) * pageSize, pageSize
	}
	return 0, DefaultPageSize
}
