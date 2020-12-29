package diff

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hero1s/golib/tools/qbtool/cmd/base"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Schema struct {
	CatalogName             string
	SchemaName              string
	DefaultCharacterSetName string
	DefaultCollationName    string
	SqlPath                 sql.NullString
}

type Table struct {
	TableCatalog   string
	TableSchema    string
	TableName      string
	TableType      string
	ENGINE         sql.NullString
	VERSION        sql.NullInt64
	RowFormat      sql.NullString
	TableRows      sql.NullInt64
	AvgRowLength   sql.NullInt64
	DataLength     sql.NullInt64
	MaxDataLength  sql.NullInt64
	IndexLength    sql.NullInt64
	DataFree       sql.NullInt64
	AutoIncrement  sql.NullInt64
	CreateTime     sql.NullTime
	UpdateTime     sql.NullTime
	CheckTime      sql.NullTime
	TableCollation sql.NullString
	CHECKSUM       sql.NullInt64
	CreateOptions  sql.NullString
	TableComment   string
}

type Column struct {
	TableCatalog           string
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int
	ColumnDefault          sql.NullString
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	CharacterOctetLength   sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	DatetimePrecision      sql.NullInt64
	CharacterSetName       sql.NullString
	CollationName          sql.NullString
	ColumnType             string
	ColumnKey              string
	EXTRA                  string
	PRIVILEGES             string
	ColumnComment          string
}

type Statistic struct {
	TableCatalog string
	TableSchema  string
	TableName    string
	NonUnique    int64
	IndexSchema  string
	IndexName    string
	SeqInIndex   int
	ColumnName   string
	COLLATION    sql.NullString
	CARDINALITY  sql.NullInt64
	SubPart      sql.NullInt32
	PACKED       sql.NullString
	NULLABLE     string
	IndexType    string
	COMMENT      sql.NullString
	IndexComment string
	IsVisible    sql.NullString
}

func RunDiff(source, target, path string) error {
	if source == "" || target == "" {
		fmt.Printf("please see %s help db\n", base.ToolName)
		return nil
	}
	if path == "" || path == "." {
		path = "diff.sql"
	}

	dir, _ := filepath.Split(path)
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	logFile, err := os.Create(path)
	defer logFile.Close()
	if err != nil {
		return err
	}

	LogPrintln := func(str string) {
		_, _ = logFile.WriteString(str + "\n")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	// 创建一个日志对象
	LogPrintln("-- ------------- 数据库差异对比----------------")
	LogPrintln("-- 对比时间: " + now)
	LogPrintln("-- 对比源数据库: " + source)
	LogPrintln("-- 对比目的数据库: " + target)
	LogPrintln("-- ------------- 数据库差异对比----------------")

	// username:password@tcp(ip:port)/db_name?charset=ut8
	replace := strings.NewReplacer("tcp(", "", ")/", "/")
	source = replace.Replace(source)
	target = replace.Replace(target)
	// username:password@ip:port/db_name?charset=ut8
	sourceDB := strings.Split(source[strings.LastIndex(source, "/")+1:], "?")[0]
	targetDB := strings.Split(target[strings.LastIndex(target, "/")+1:], "?")[0]
	var sourceUser = strings.Split(source[0:strings.LastIndex(source, "@")], ":")
	var sourceHost = strings.Split(strings.Split(source[strings.LastIndex(source, "@")+1:], "/")[0], ":")

	var targetUser = strings.Split(target[0:strings.LastIndex(target, "@")], ":")
	var targetHost = strings.Split(strings.Split(target[strings.LastIndex(target, "@")+1:], "/")[0], ":")

	sourceDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?parseTime=true&charset=utf8", sourceUser[0], sourceUser[1], sourceHost[0], sourceHost[1]))
	defer sourceDb.Close()
	if err != nil {
		return err
	}

	targetDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?parseTime=true&charset=utf8", targetUser[0], targetUser[1], targetHost[0], targetHost[1]))
	defer targetDb.Close()
	if err != nil {
		return err
	}
	var sourceSchemaCount int32
	var targetSchemaCount int32

	err = sourceDb.QueryRow("SELECT COUNT(*) FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` = ?", sourceDB).Scan(&sourceSchemaCount)

	if err != nil {
		return err
	}

	if sourceSchemaCount <= 0 {
		return fmt.Errorf("源数据库 `%s` 不存在。", sourceDB)
	}

	err = targetDb.QueryRow("SELECT COUNT(*) FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` = ?", targetDB).Scan(&targetSchemaCount)

	if err != nil {
		return err
	}

	if targetSchemaCount <= 0 {
		return fmt.Errorf("目标数据库 `%s` 不存在。", targetDB)
	}

	var sourceSchema Schema
	var targetSchema Schema

	err = sourceDb.QueryRow("SELECT `CATALOG_NAME`, `SCHEMA_NAME`, `DEFAULT_CHARACTER_SET_NAME`, `DEFAULT_COLLATION_NAME`, `SQL_PATH` "+
		"FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` = ?", sourceDB).Scan(
		&sourceSchema.CatalogName,
		&sourceSchema.SchemaName,
		&sourceSchema.DefaultCharacterSetName,
		&sourceSchema.DefaultCollationName,
		&sourceSchema.SqlPath,
	)

	if err != nil {
		return err
	}

	err = targetDb.QueryRow("SELECT `CATALOG_NAME`, `SCHEMA_NAME`, `DEFAULT_CHARACTER_SET_NAME`, `DEFAULT_COLLATION_NAME`, `SQL_PATH` "+
		"FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` = ?", targetDB).Scan(
		&targetSchema.CatalogName,
		&targetSchema.SchemaName,
		&targetSchema.DefaultCharacterSetName,
		&targetSchema.DefaultCollationName,
		&targetSchema.SqlPath,
	)

	if err != nil {
		return err
	}

	sourceTableRows, err := sourceDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `TABLE_TYPE`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `AVG_ROW_LENGTH`, `DATA_LENGTH`, `MAX_DATA_LENGTH`, `INDEX_LENGTH`, `DATA_FREE`, `AUTO_INCREMENT`, `CREATE_TIME`, `UPDATE_TIME`, `CHECK_TIME`, `TABLE_COLLATION`, `CHECKSUM`, `CREATE_OPTIONS`, `TABLE_COMMENT` "+
		"FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` = ? ORDER BY `TABLE_NAME` ASC", sourceDB)

	if err != nil {
		return err
	}

	var sourceTableData []Table
	sourceTableMap := make(map[string]Table)

	for sourceTableRows.Next() {
		var table Table

		err := sourceTableRows.Scan(
			&table.TableCatalog,
			&table.TableSchema,
			&table.TableName,
			&table.TableType,
			&table.ENGINE,
			&table.VERSION,
			&table.RowFormat,
			&table.TableRows,
			&table.AvgRowLength,
			&table.DataLength,
			&table.MaxDataLength,
			&table.IndexLength,
			&table.DataFree,
			&table.AutoIncrement,
			&table.CreateTime,
			&table.UpdateTime,
			&table.CheckTime,
			&table.TableCollation,
			&table.CHECKSUM,
			&table.CreateOptions,
			&table.TableComment,
		)

		if err != nil {
			return err
		}

		sourceTableData = append(sourceTableData, table)
		sourceTableMap[table.TableName] = table
	}

	targetTableRows, err := targetDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `TABLE_TYPE`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `AVG_ROW_LENGTH`, `DATA_LENGTH`, `MAX_DATA_LENGTH`, `INDEX_LENGTH`, `DATA_FREE`, `AUTO_INCREMENT`, `CREATE_TIME`, `UPDATE_TIME`, `CHECK_TIME`, `TABLE_COLLATION`, `CHECKSUM`, `CREATE_OPTIONS`, `TABLE_COMMENT` "+
		"FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` = ? ORDER BY `TABLE_NAME` ASC", targetDB)

	if err != nil {
		return err
	}

	var targetTableData []Table
	targetTableMap := make(map[string]Table)

	for targetTableRows.Next() {
		var table Table

		err := targetTableRows.Scan(
			&table.TableCatalog,
			&table.TableSchema,
			&table.TableName,
			&table.TableType,
			&table.ENGINE,
			&table.VERSION,
			&table.RowFormat,
			&table.TableRows,
			&table.AvgRowLength,
			&table.DataLength,
			&table.MaxDataLength,
			&table.IndexLength,
			&table.DataFree,
			&table.AutoIncrement,
			&table.CreateTime,
			&table.UpdateTime,
			&table.CheckTime,
			&table.TableCollation,
			&table.CHECKSUM,
			&table.CreateOptions,
			&table.TableComment,
		)

		if err != nil {
			return err
		}

		targetTableData = append(targetTableData, table)
		targetTableMap[table.TableName] = table
	}

	// 获取源数据库collation跟character set的对应关系
	sourceCharset, err := sourceDb.Query("SELECT `COLLATION_NAME`, `CHARACTER_SET_NAME`" +
		"FROM `information_schema`.`COLLATION_CHARACTER_SET_APPLICABILITY`")

	if err != nil {
		return err
	}
	var sourceCharsetMap = make(map[string]string)
	for sourceCharset.Next() {
		var charset, collation string
		err := sourceCharset.Scan(&collation, &charset)
		if err != nil {
			return err
		}
		sourceCharsetMap[collation] = charset
	}

	var diffSql []string

	// DROP TABLE...
	for _, targetTable := range targetTableData {
		if _, ok := sourceTableMap[targetTable.TableName]; !ok {
			diffSql = append(diffSql, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", targetTable.TableName))
		}
	}

	// ALTER TABLE...
	for _, sourceTable := range sourceTableData {
		if targetTable, ok := targetTableMap[sourceTable.TableName]; ok {
			if !CompareTable(sourceTable, targetTable) {
				var (
					engine           = "" // ENGINE = MEMORY, ROW_FORMAT = Fixed  // ROW_FORMAT需要放在最后面
					charset          = "" // CHARACTER SET = utf16, COLLATE = utf16_german2_ci
					comment          = "" // COMMENT = '账号表'
					alterTable       = fmt.Sprintf("ALTER TABLE `%s`", sourceTable.TableName)
					originAlterTable = alterTable
					engineChangeFlag = false
				)
				if sourceTable.ENGINE.Valid && sourceTable.ENGINE.String != targetTable.ENGINE.String {
					engineChangeFlag = true
					engine = fmt.Sprintf(" ENGINE=%s", sourceTable.ENGINE.String)
					alterTable = alterTable + engine
				}
				if sourceTable.TableCollation.Valid && sourceTable.TableCollation.String != targetTable.TableCollation.String {
					collation := sourceTable.TableCollation.String
					character, ok := sourceCharsetMap[collation]
					if !ok {
						return errors.New(fmt.Sprintf("源(source)数据库找不到COLLATION:%s对应的CHARACTER", collation))
					}
					charset = fmt.Sprintf(" CHARACTER SET=%s, COLLATE=%s", character, collation)
					if alterTable == originAlterTable {
						alterTable = alterTable + charset
					} else {
						alterTable = alterTable + "," + charset
					}
				}
				if sourceTable.TableComment != targetTable.TableComment {
					comment = fmt.Sprintf(" COMMENT='%s'", sourceTable.TableComment)
					if alterTable == originAlterTable {
						alterTable = alterTable + comment
					} else {
						alterTable = alterTable + "," + comment
					}
				}
				if engineChangeFlag {
					alterTable = alterTable + fmt.Sprintf(",ROW_FORMAT=%s", sourceTable.RowFormat.String)
				}
				alterTable = alterTable + ";"
				diffSql = append(diffSql, alterTable)
			}
		}

	}

	for _, sourceTable := range sourceTableData {
		if _, ok := targetTableMap[sourceTable.TableName]; ok {
			// ALTER TABLE ...
			sourceColumnRows, err := sourceDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `ORDINAL_POSITION`, `COLUMN_DEFAULT`, `IS_NULLABLE`, `DATA_TYPE`, `CHARACTER_MAXIMUM_LENGTH`, `CHARACTER_OCTET_LENGTH`, `NUMERIC_PRECISION`, `NUMERIC_SCALE`, `DATETIME_PRECISION`, `CHARACTER_SET_NAME`, `COLLATION_NAME`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`, `PRIVILEGES`, `COLUMN_COMMENT` "+
				"FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY `ORDINAL_POSITION` ASC", sourceDB, sourceTable.TableName)

			if err != nil {
				return err
			}

			var sourceColumnData []Column

			for sourceColumnRows.Next() {
				var column Column

				err := sourceColumnRows.Scan(
					&column.TableCatalog,
					&column.TableSchema,
					&column.TableName,
					&column.ColumnName,
					&column.OrdinalPosition,
					&column.ColumnDefault,
					&column.IsNullable,
					&column.DataType,
					&column.CharacterMaximumLength,
					&column.CharacterOctetLength,
					&column.NumericPrecision,
					&column.NumericScale,
					&column.DatetimePrecision,
					&column.CharacterSetName,
					&column.CollationName,
					&column.ColumnType,
					&column.ColumnKey,
					&column.EXTRA,
					&column.PRIVILEGES,
					&column.ColumnComment,
				)

				if err != nil {
					return err
				}

				sourceColumnData = append(sourceColumnData, column)
			}

			sourceColumnDataLen := len(sourceColumnData)

			targetColumnRows, err := targetDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `ORDINAL_POSITION`, `COLUMN_DEFAULT`, `IS_NULLABLE`, `DATA_TYPE`, `CHARACTER_MAXIMUM_LENGTH`, `CHARACTER_OCTET_LENGTH`, `NUMERIC_PRECISION`, `NUMERIC_SCALE`, `DATETIME_PRECISION`, `CHARACTER_SET_NAME`, `COLLATION_NAME`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`, `PRIVILEGES`, `COLUMN_COMMENT` "+
				"FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY `ORDINAL_POSITION` ASC", targetDB, sourceTable.TableName)

			if err != nil {
				return err
			}

			var targetColumnData []Column

			for targetColumnRows.Next() {
				var column Column

				err := targetColumnRows.Scan(
					&column.TableCatalog,
					&column.TableSchema,
					&column.TableName,
					&column.ColumnName,
					&column.OrdinalPosition,
					&column.ColumnDefault,
					&column.IsNullable,
					&column.DataType,
					&column.CharacterMaximumLength,
					&column.CharacterOctetLength,
					&column.NumericPrecision,
					&column.NumericScale,
					&column.DatetimePrecision,
					&column.CharacterSetName,
					&column.CollationName,
					&column.ColumnType,
					&column.ColumnKey,
					&column.EXTRA,
					&column.PRIVILEGES,
					&column.ColumnComment,
				)

				if err != nil {
					return err
				}

				targetColumnData = append(targetColumnData, column)
			}

			targetColumnDataLen := len(targetColumnData)

			// ALTER LIST ...
			var (
				alterTableSql  []string
				alterColumnSql []string
				alterKeySql    []string
			)

			if sourceColumnDataLen > 0 && targetColumnDataLen > 0 {
				sourceColumns := make(map[string]Column)
				targetColumns := make(map[string]Column)
				sourceColumnsPos := make(map[int]Column)
				targetColumnsPos := make(map[int]Column)

				for _, sourceColumn := range sourceColumnData {
					sourceColumns[sourceColumn.ColumnName] = sourceColumn
					sourceColumnsPos[sourceColumn.OrdinalPosition] = sourceColumn
				}

				for _, targetColumn := range targetColumnData {
					targetColumns[targetColumn.ColumnName] = targetColumn
					targetColumnsPos[targetColumn.OrdinalPosition] = targetColumn
				}

				if !CompareColumns(sourceColumnsPos, targetColumnsPos) {
					alterTableSql = append(alterTableSql, fmt.Sprintf("ALTER TABLE `%s`", sourceTable.TableName))

					// DROP COLUMN ...
					for _, targetColumn := range targetColumns {
						if _, ok := sourceColumns[targetColumn.ColumnName]; !ok {
							ResetCalcPosition(targetColumn.ColumnName, targetColumn.OrdinalPosition, targetColumns, 3)

							alterColumnSql = append(alterColumnSql, fmt.Sprintf("  DROP COLUMN `%s`", targetColumn.ColumnName))
						}
					}

					// ADD COLUMN ...
					for _, sourceColumn := range sourceColumnData {
						columnName := sourceColumn.ColumnName

						if _, ok := targetColumns[columnName]; !ok {
							nullAbleDefault := GetColumnNullAbleDefault(sourceColumn)
							var (
								character = ""
								extra     = ""
								comment   = " COMMENT '%s'"
							)

							if sourceColumn.CharacterSetName.Valid {
								if sourceColumn.CharacterSetName.String != sourceSchema.DefaultCharacterSetName {
									character = fmt.Sprintf(" CHARACTER SET %s", sourceColumn.CharacterSetName.String)
								}
							}

							if sourceColumn.EXTRA != "" {
								extra = fmt.Sprintf(" %s", strings.ToUpper(sourceColumn.EXTRA))
							}

							after := GetColumnAfter(sourceColumn.OrdinalPosition, sourceColumnsPos)
							comment = fmt.Sprintf(comment, sourceColumn.ColumnComment)

							ResetCalcPosition(columnName, sourceColumn.OrdinalPosition, targetColumns, 1)

							alterColumnSql = append(alterColumnSql, fmt.Sprintf("  ADD COLUMN `%s` %s%s%s%s%s %s", columnName, sourceColumn.ColumnType, character, nullAbleDefault, extra, comment, after))
						}
					}

					// MODIFY COLUMN ...
					for _, sourceColumn := range sourceColumnData {
						columnName := sourceColumn.ColumnName

						if _, ok := targetColumns[columnName]; ok {
							if !CompareColumn(sourceColumn, targetColumns[columnName]) {
								nullAbleDefault := GetColumnNullAbleDefault(sourceColumn)

								var (
									character = ""
									extra     = ""
									comment   = " COMMENT '%s'"
								)

								if sourceColumn.CharacterSetName.Valid {
									if sourceColumn.CharacterSetName.String != sourceSchema.DefaultCharacterSetName {
										character = fmt.Sprintf(" CHARACTER SET %s", sourceColumn.CharacterSetName.String)
									}
								}

								if sourceColumn.EXTRA != "" {
									extra = fmt.Sprintf(" %s", strings.ToUpper(sourceColumn.EXTRA))
								}

								after := GetColumnAfter(sourceColumn.OrdinalPosition, sourceColumnsPos)
								comment = fmt.Sprintf(comment, sourceColumn.ColumnComment)

								ResetCalcPosition(columnName, sourceColumn.OrdinalPosition, targetColumns, 2)

								alterColumnSql = append(alterColumnSql, fmt.Sprintf("  MODIFY COLUMN `%s` %s%s%s%s%s %s", columnName, sourceColumn.ColumnType, character, nullAbleDefault, extra, comment, after))
							}
						}
					}
				}
			}

			// ADD KEY AND DROP INDEX ...
			sourceStatisticsRows, err := sourceDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `NON_UNIQUE`, `INDEX_SCHEMA`, `INDEX_NAME`, `SEQ_IN_INDEX`, `COLUMN_NAME`, `COLLATION`, `CARDINALITY`, `SUB_PART`, `PACKED`, `NULLABLE`, `INDEX_TYPE`, `COMMENT`, `INDEX_COMMENT` "+
				"FROM `information_schema`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?", sourceDB, sourceTable.TableName)

			if err != nil {
				return err
			}

			var sourceStatisticsData []Statistic

			for sourceStatisticsRows.Next() {
				var statistic Statistic

				err := sourceStatisticsRows.Scan(
					&statistic.TableCatalog,
					&statistic.TableSchema,
					&statistic.TableName,
					&statistic.NonUnique,
					&statistic.IndexSchema,
					&statistic.IndexName,
					&statistic.SeqInIndex,
					&statistic.ColumnName,
					&statistic.COLLATION,
					&statistic.CARDINALITY,
					&statistic.SubPart,
					&statistic.PACKED,
					&statistic.NULLABLE,
					&statistic.IndexType,
					&statistic.COMMENT,
					&statistic.IndexComment,
				)

				if err != nil {
					return err
				}

				sourceStatisticsData = append(sourceStatisticsData, statistic)
			}

			targetStatisticsRows, err := targetDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `NON_UNIQUE`, `INDEX_SCHEMA`, `INDEX_NAME`, `SEQ_IN_INDEX`, `COLUMN_NAME`, `COLLATION`, `CARDINALITY`, `SUB_PART`, `PACKED`, `NULLABLE`, `INDEX_TYPE`, `COMMENT`, `INDEX_COMMENT` "+
				"FROM `information_schema`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?", targetDB, sourceTable.TableName)

			if err != nil {
				return err
			}

			var targetStatisticsData []Statistic

			for targetStatisticsRows.Next() {
				var statistic Statistic

				err := targetStatisticsRows.Scan(
					&statistic.TableCatalog,
					&statistic.TableSchema,
					&statistic.TableName,
					&statistic.NonUnique,
					&statistic.IndexSchema,
					&statistic.IndexName,
					&statistic.SeqInIndex,
					&statistic.ColumnName,
					&statistic.COLLATION,
					&statistic.CARDINALITY,
					&statistic.SubPart,
					&statistic.PACKED,
					&statistic.NULLABLE,
					&statistic.IndexType,
					&statistic.COMMENT,
					&statistic.IndexComment,
				)

				if err != nil {
					return err
				}

				targetStatisticsData = append(targetStatisticsData, statistic)
			}

			sourceStatisticsDataLen := len(sourceStatisticsData)

			if sourceStatisticsDataLen > 0 {
				sourceStatisticsDataMap := make(map[string]map[int]Statistic)
				targetStatisticsDataMap := make(map[string]map[int]Statistic)

				for _, sourceStatistic := range sourceStatisticsData {
					if _, ok := sourceStatisticsDataMap[sourceStatistic.IndexName]; ok {
						sourceStatisticsDataMap[sourceStatistic.IndexName][sourceStatistic.SeqInIndex] = sourceStatistic
					} else {
						sourceSeqInIndexStatisticMap := make(map[int]Statistic)
						sourceSeqInIndexStatisticMap[sourceStatistic.SeqInIndex] = sourceStatistic
						sourceStatisticsDataMap[sourceStatistic.IndexName] = sourceSeqInIndexStatisticMap
					}
				}

				for _, targetStatistic := range targetStatisticsData {
					if _, ok := targetStatisticsDataMap[targetStatistic.IndexName]; ok {
						targetStatisticsDataMap[targetStatistic.IndexName][targetStatistic.SeqInIndex] = targetStatistic
					} else {
						targetSeqInIndexStatisticMap := make(map[int]Statistic)
						targetSeqInIndexStatisticMap[targetStatistic.SeqInIndex] = targetStatistic
						targetStatisticsDataMap[targetStatistic.IndexName] = targetSeqInIndexStatisticMap
					}
				}

				if !CompareStatistics(sourceStatisticsDataMap, targetStatisticsDataMap) {
					if len(alterTableSql) <= 0 {
						alterTableSql = append(alterTableSql, fmt.Sprintf("ALTER TABLE `%s`", sourceTable.TableName))
					}

					// DROP INDEX ...
					for targetIndexName, _ := range targetStatisticsDataMap {
						if _, ok := sourceStatisticsDataMap[targetIndexName]; !ok {
							if "PRIMARY" == targetIndexName {
								alterKeySql = append(alterKeySql, "  DROP PRIMARY KEY")
							} else {
								alterKeySql = append(alterKeySql, fmt.Sprintf("  DROP INDEX `%s`", targetIndexName))
							}
						}
					}

					// DROP INDEX ... AND ADD KEY ...
					for sourceIndexName, sourceStatisticMap := range sourceStatisticsDataMap {
						if _, ok := targetStatisticsDataMap[sourceIndexName]; ok {
							if !CompareStatisticsIndex(sourceStatisticMap, targetStatisticsDataMap[sourceIndexName]) {
								// DROP INDEX ...
								if "PRIMARY" == sourceIndexName {
									alterKeySql = append(alterKeySql, "  DROP PRIMARY KEY")
								} else {
									alterKeySql = append(alterKeySql, fmt.Sprintf("  DROP INDEX `%s`", sourceIndexName))
								}

								// ADD KEY ...
								alterKeySql = append(alterKeySql, fmt.Sprintf("  ADD %s", GetAddKeys(sourceIndexName, sourceStatisticMap)))
							}
						} else {
							// ADD KEY ...
							alterKeySql = append(alterKeySql, fmt.Sprintf("  ADD %s", GetAddKeys(sourceIndexName, sourceStatisticMap)))
						}
					}

					if len(alterKeySql) > 0 {
						for _, keySql := range alterKeySql {
							alterColumnSql = append(alterColumnSql, keySql)
						}
					}
				}
			}

			// ALTER TABLE SQL ...
			alterColumnSqlLen := len(alterColumnSql)

			if alterColumnSqlLen > 0 {
				for _, alterColumn := range alterColumnSql {
					var columnDot = ""
					if alterColumn == alterColumnSql[alterColumnSqlLen-1] {
						columnDot = ";"
					} else {
						columnDot = ","
					}

					alterTableSql = append(alterTableSql, fmt.Sprintf("%s%s", alterColumn, columnDot))
				}
			}

			alterTableSqlLen := len(alterTableSql)

			if alterTableSqlLen > 0 {
				diffSql = append(diffSql, strings.Join(alterTableSql, "\n"))
			}
		} else {
			// CREATE TABLE ...
			sourceColumnRows, err := sourceDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `ORDINAL_POSITION`, `COLUMN_DEFAULT`, `IS_NULLABLE`, `DATA_TYPE`, `CHARACTER_MAXIMUM_LENGTH`, `CHARACTER_OCTET_LENGTH`, `NUMERIC_PRECISION`, `NUMERIC_SCALE`, `DATETIME_PRECISION`, `CHARACTER_SET_NAME`, `COLLATION_NAME`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`, `PRIVILEGES`, `COLUMN_COMMENT` "+
				"FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY `ORDINAL_POSITION` ASC", sourceDB, sourceTable.TableName)

			if err != nil {
				return err
			}

			var sourceColumnData []Column

			for sourceColumnRows.Next() {
				var column Column

				err = sourceColumnRows.Scan(
					&column.TableCatalog,
					&column.TableSchema,
					&column.TableName,
					&column.ColumnName,
					&column.OrdinalPosition,
					&column.ColumnDefault,
					&column.IsNullable,
					&column.DataType,
					&column.CharacterMaximumLength,
					&column.CharacterOctetLength,
					&column.NumericPrecision,
					&column.NumericScale,
					&column.DatetimePrecision,
					&column.CharacterSetName,
					&column.CollationName,
					&column.ColumnType,
					&column.ColumnKey,
					&column.EXTRA,
					&column.PRIVILEGES,
					&column.ColumnComment,
				)

				sourceColumnData = append(sourceColumnData, column)
			}

			if err != nil {
				return err
			}

			sourceColumnDataLen := len(sourceColumnData)

			if sourceColumnDataLen > 0 {
				sourceStatisticsRows, err := sourceDb.Query("SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `NON_UNIQUE`, `INDEX_SCHEMA`, `INDEX_NAME`, `SEQ_IN_INDEX`, `COLUMN_NAME`, `COLLATION`, `CARDINALITY`, `SUB_PART`, `PACKED`, `NULLABLE`, `INDEX_TYPE`, `COMMENT`, `INDEX_COMMENT` "+
					"FROM `information_schema`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?", sourceDB, sourceTable.TableName)

				if err != nil {
					return err
				}

				var sourceStatisticsData []Statistic

				for sourceStatisticsRows.Next() {
					var statistic Statistic

					err := sourceStatisticsRows.Scan(
						&statistic.TableCatalog,
						&statistic.TableSchema,
						&statistic.TableName,
						&statistic.NonUnique,
						&statistic.IndexSchema,
						&statistic.IndexName,
						&statistic.SeqInIndex,
						&statistic.ColumnName,
						&statistic.COLLATION,
						&statistic.CARDINALITY,
						&statistic.SubPart,
						&statistic.PACKED,
						&statistic.NULLABLE,
						&statistic.IndexType,
						&statistic.COMMENT,
						&statistic.IndexComment,
					)

					if err != nil {
						return err
					}

					sourceStatisticsData = append(sourceStatisticsData, statistic)
				}

				var createTableSql []string

				createTableSql = append(createTableSql, fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (", sourceTable.TableName))

				//fmt.Println(sourceTable.TableName)

				// COLUMNS ...
				for _, sourceColumn := range sourceColumnData {
					var (
						character = ""
						extra     = ""
						comment   = " COMMENT '%s'"
						dot       = ""
					)

					if sourceColumn.CharacterSetName.Valid {
						if sourceColumn.CharacterSetName.String != sourceSchema.DefaultCharacterSetName {
							character = fmt.Sprintf(" CHARACTER SET %s", sourceColumn.CharacterSetName.String)
						}
					}

					if sourceColumn.EXTRA != "" {
						extra = fmt.Sprintf(" %s", strings.ToUpper(sourceColumn.EXTRA))
					}

					comment = fmt.Sprintf(comment, sourceColumn.ColumnComment)

					if sourceColumn != sourceColumnData[sourceColumnDataLen-1] || len(sourceStatisticsData) > 0 {
						dot = ","
					}

					createTableSql = append(createTableSql, fmt.Sprintf("  `%s` %s%s%s%s%s%s", sourceColumn.ColumnName, sourceColumn.ColumnType, character, GetColumnNullAbleDefault(sourceColumn), extra, comment, dot))

					//fmt.Println(sourceColumn.ColumnName)
				}

				// KEY ...
				var createKeySql []string
				sourceStatisticsLen := len(sourceStatisticsData)

				if sourceStatisticsLen > 0 {
					var sourceStatisticIndexNameArray []string
					sourceStatisticsDataMap := make(map[string]map[int]Statistic)

					for _, sourceStatistic := range sourceStatisticsData {
						if _, ok := sourceStatisticsDataMap[sourceStatistic.IndexName]; ok {
							sourceStatisticsDataMap[sourceStatistic.IndexName][sourceStatistic.SeqInIndex] = sourceStatistic
						} else {
							sourceSeqInIndexStatisticMap := make(map[int]Statistic)
							sourceSeqInIndexStatisticMap[sourceStatistic.SeqInIndex] = sourceStatistic
							sourceStatisticsDataMap[sourceStatistic.IndexName] = sourceSeqInIndexStatisticMap
						}

						if !inArray(sourceStatistic.IndexName, sourceStatisticIndexNameArray) {
							sourceStatisticIndexNameArray = append(sourceStatisticIndexNameArray, sourceStatistic.IndexName)
						}
					}

					for _, sourceIndexName := range sourceStatisticIndexNameArray {
						createKeySql = append(createKeySql, fmt.Sprintf("  %s", GetAddKeys(sourceIndexName, sourceStatisticsDataMap[sourceIndexName])))
					}
				}

				createTableSql = append(createTableSql, strings.Join(createKeySql, ",\n"))
				createTableSql = append(createTableSql, fmt.Sprintf(") ENGINE=%s DEFAULT CHARSET=%s;", sourceTable.ENGINE.String, sourceSchema.DefaultCharacterSetName))

				diffSql = append(diffSql, strings.Join(createTableSql, "\n"))
			}
		}
	}

	// Print Sql...
	if len(diffSql) > 0 {
		LogPrintln(fmt.Sprintf("SET NAMES %s;\n", sourceSchema.DefaultCharacterSetName))
		LogPrintln(strings.Join(diffSql, "\n\n"))
	}

	return nil

}

func GetColumnNullAbleDefault(column Column) string {
	var nullAbleDefault = ""

	if column.IsNullable == "NO" {
		if column.ColumnDefault.Valid {
			if column.ColumnDefault.String == "CURRENT_TIMESTAMP" {
				nullAbleDefault = fmt.Sprintf(" NOT NULL DEFAULT %s", column.ColumnDefault.String)
			} else {
				nullAbleDefault = fmt.Sprintf(" NOT NULL DEFAULT '%s'", column.ColumnDefault.String)
			}
		} else {
			nullAbleDefault = " NOT NULL"
		}
	} else {
		if column.ColumnDefault.Valid {
			if column.ColumnDefault.String == "CURRENT_TIMESTAMP" {
				nullAbleDefault = fmt.Sprintf(" NULL DEFAULT %s", column.ColumnDefault.String)
			} else {
				nullAbleDefault = fmt.Sprintf(" DEFAULT '%s'", column.ColumnDefault.String)
			}
		} else {
			nullAbleDefault = " DEFAULT NULL"
		}
	}

	return nullAbleDefault
}

func GetAddKeys(indexName string, statisticMap map[int]Statistic) string {
	if 1 == statisticMap[1].NonUnique {
		var seqInIndexSort []int
		var columnNames []string

		for seqInIndex, _ := range statisticMap {
			seqInIndexSort = append(seqInIndexSort, seqInIndex)
		}

		sort.Ints(seqInIndexSort)

		for _, seqInIndex := range seqInIndexSort {
			var subPart = ""

			if statisticMap[seqInIndex].SubPart.Valid {
				subPart = fmt.Sprintf("(%d)", statisticMap[seqInIndex].SubPart.Int32)
			}

			columnNames = append(columnNames, fmt.Sprintf("`%s`%s", statisticMap[seqInIndex].ColumnName, subPart))
		}

		return fmt.Sprintf("KEY `%s` (%s)", indexName, strings.Join(columnNames, ","))
	} else {
		if "PRIMARY" == indexName {
			var seqInIndexSort []int
			var columnNames []string

			for seqInIndex, _ := range statisticMap {
				seqInIndexSort = append(seqInIndexSort, seqInIndex)
			}

			sort.Ints(seqInIndexSort)

			for _, seqInIndex := range seqInIndexSort {
				var subPart = ""

				if statisticMap[seqInIndex].SubPart.Valid {
					subPart = fmt.Sprintf("(%d)", statisticMap[seqInIndex].SubPart.Int32)
				}

				columnNames = append(columnNames, fmt.Sprintf("`%s`%s", statisticMap[seqInIndex].ColumnName, subPart))
			}

			return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(columnNames, ","))
		} else {
			var seqInIndexSort []int
			var columnNames []string

			for seqInIndex, _ := range statisticMap {
				seqInIndexSort = append(seqInIndexSort, seqInIndex)
			}

			sort.Ints(seqInIndexSort)

			for _, seqInIndex := range seqInIndexSort {
				var subPart = ""

				if statisticMap[seqInIndex].SubPart.Valid {
					subPart = fmt.Sprintf("(%d)", statisticMap[seqInIndex].SubPart.Int32)
				}

				columnNames = append(columnNames, fmt.Sprintf("`%s`%s", statisticMap[seqInIndex].ColumnName, subPart))
			}

			return fmt.Sprintf("UNIQUE KEY `%s` (%s)", indexName, strings.Join(columnNames, ","))
		}
	}
}

func CompareColumns(sourceColumnsPos map[int]Column, targetColumnsPos map[int]Column) bool {
	if len(sourceColumnsPos) != len(targetColumnsPos) {
		return false
	} else {
		for sourcePos, sourceColumn := range sourceColumnsPos {
			if _, ok := targetColumnsPos[sourcePos]; ok {
				targetColumn := targetColumnsPos[sourcePos]

				if !CompareColumn(sourceColumn, targetColumn) {
					return false
				}
			} else {
				return false
			}

		}
	}

	return true
}

// 比对表属性
func CompareTable(sourceTable, targetTable Table) bool {
	if sourceTable.ENGINE != targetTable.ENGINE {
		return false
	}
	if sourceTable.RowFormat != targetTable.RowFormat {
		return false
	}
	if sourceTable.TableCollation != targetTable.TableCollation {
		return false
	}
	if sourceTable.TableComment != targetTable.TableComment {
		return false
	}
	return true
}

func CompareColumn(sourceColumn Column, targetColumn Column) bool {
	if sourceColumn.ColumnName != targetColumn.ColumnName {
		return false
	}

	if sourceColumn.OrdinalPosition != targetColumn.OrdinalPosition {
		return false
	}

	if sourceColumn.ColumnDefault != targetColumn.ColumnDefault {
		return false
	}

	if sourceColumn.IsNullable != targetColumn.IsNullable {
		return false
	}

	if sourceColumn.DataType != targetColumn.DataType {
		return false
	}

	if sourceColumn.CharacterMaximumLength != targetColumn.CharacterMaximumLength {
		return false
	}

	//禁用实际精度检验，因为 TiDB 和 MySQL 在设置不标准的情况下，值会不一样。
	//if sourceColumn.NumericPrecision != targetColumn.NumericPrecision {
	//	return false
	//}

	if sourceColumn.NumericScale != targetColumn.NumericScale {
		return false
	}

	if sourceColumn.DatetimePrecision != targetColumn.DatetimePrecision {
		return false
	}

	if sourceColumn.CharacterSetName != targetColumn.CharacterSetName {
		return false
	}

	if sourceColumn.CollationName != targetColumn.CollationName {
		return false
	}

	if sourceColumn.ColumnType != targetColumn.ColumnType {
		return false
	}

	if sourceColumn.EXTRA != targetColumn.EXTRA {
		return false
	}
	if sourceColumn.ColumnComment != targetColumn.ColumnComment {
		return false
	}

	return true
}

func CompareStatistics(sourceStatisticsMap map[string]map[int]Statistic, targetStatisticsMap map[string]map[int]Statistic) bool {
	if len(sourceStatisticsMap) != len(targetStatisticsMap) {
		return false
	} else {
		for indexName, sourceStatisticMap := range sourceStatisticsMap {
			if _, ok := targetStatisticsMap[indexName]; ok {
				if len(sourceStatisticMap) != len(targetStatisticsMap[indexName]) {
					return false
				} else {
					for seqInIndex, sourceStatistic := range sourceStatisticMap {
						if _, ok := targetStatisticsMap[indexName][seqInIndex]; ok {
							if !CompareStatistic(sourceStatistic, targetStatisticsMap[indexName][seqInIndex]) {
								return false
							}
						} else {
							return false
						}
					}
				}
			} else {
				return false
			}
		}
	}

	return true
}

func CompareStatisticsIndex(sourceStatisticMap map[int]Statistic, targetStatisticMap map[int]Statistic) bool {
	if len(sourceStatisticMap) != len(targetStatisticMap) {
		return false
	} else {
		for seqInIndex, sourceStatistic := range sourceStatisticMap {
			if _, ok := targetStatisticMap[seqInIndex]; ok {
				if !CompareStatistic(sourceStatistic, targetStatisticMap[seqInIndex]) {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

func CompareStatistic(sourceStatistic Statistic, targetStatistic Statistic) bool {
	if sourceStatistic.NonUnique != targetStatistic.NonUnique {
		return false
	}

	if sourceStatistic.IndexName != targetStatistic.IndexName {
		return false
	}

	if sourceStatistic.SeqInIndex != targetStatistic.SeqInIndex {
		return false
	}

	if sourceStatistic.ColumnName != targetStatistic.ColumnName {
		return false
	}

	if sourceStatistic.SubPart != targetStatistic.SubPart {
		return false
	}

	if sourceStatistic.IndexType != targetStatistic.IndexType {
		return false
	}

	return true
}

func ResetCalcPosition(columnName string, sourcePos int, targetColumns map[string]Column, status int) {
	switch status {
	case 1:
		// ADD ...
		for targetColumnName, targetColumn := range targetColumns {
			if targetColumn.OrdinalPosition >= sourcePos {
				targetColumn.OrdinalPosition += 1

				targetColumns[targetColumnName] = targetColumn
			}
		}
		break
	case 2:
		// MODIFY ...
		if _, ok := targetColumns[columnName]; ok {
			targetColumn := targetColumns[columnName]

			targetColumn.OrdinalPosition = sourcePos

			targetColumns[columnName] = targetColumn
		}
		break
	case 3:
		// DROP ...
		for targetColumnName, targetColumn := range targetColumns {
			if targetColumn.OrdinalPosition >= sourcePos {
				targetColumn.OrdinalPosition -= 1

				targetColumns[targetColumnName] = targetColumn
			}
		}
		break
	}
}

func GetColumnAfter(ordinalPosition int, columnsPos map[int]Column) string {
	pos := ordinalPosition - 1

	if _, ok := columnsPos[pos]; ok {
		return fmt.Sprintf("AFTER `%s`", columnsPos[pos].ColumnName)
	} else {
		return "FIRST"
	}
}

func inArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
