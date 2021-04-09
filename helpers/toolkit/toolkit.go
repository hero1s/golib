package toolkit

import (
	"fmt"
	"strconv"
)

//拼接sql id集
func MakeIdsSqlIn(ids []int) string {
	if len(ids) == 0 {
		ids = append(ids, 0)
	}
	sqlEx := ""
	tmpFlag := false
	for _, v := range ids {
		if tmpFlag {
			sqlEx += fmt.Sprintf(`,%v`, v)
		} else {
			sqlEx += fmt.Sprintf(`%v`, v)
			tmpFlag = true
		}
	}
	return sqlEx
}

// 强制转为浮点数
func Force2Float(value interface{}) float64 {
	v, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return 0
	}
	return v
}

// 强制转为整数，如果是浮点数强转为整数，会有精度损失
func Force2Int(value interface{}) int {
	return int(Force2Float(value))
}

// 强制转为整数，如果是浮点数强转为整数，会有精度损失
func Force2Int64(value interface{}) int64 {
	return int64(Force2Float(value))
}
