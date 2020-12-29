package reverse

import "github.com/hero1s/golib/tools/qbtool/cmd/base"
import "fmt"

func RunReverse(source, path string,singleFile bool) error {
	if source == "" {
		fmt.Printf("please see %s help db\n", base.ToolName)
		return nil
	}

	dbStruct := NewDBStruct()
	return dbStruct.
		Dsn(source).
		StructNameFmt(FmtUnderlineToStartUpHump).
		FieldNameFmt(FmtUnderlineToStartUpHump).
		FileNameFmt(FmtUnderline).
		SingleFile(singleFile).
		GenTableName("TableName").
		GenFullTableName("FullTableName").
		GenDbName("DbName").
		TagJson(true).
		ModelPath(path).
		//TagOrm(true).
		//AppendTag(NewTag("xml", FmtDefault)).
		Generate()
}
