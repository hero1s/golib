package datetime

import (
	"strconv"
	"time"
)

//根据时间戳,获取星座
func Constellation(tt int64) string {
	t := time.Unix(tt, 0).Format("0102")
	d, _ := strconv.ParseInt(t, 10, 64)
	if d >= 321 && d <= 419 {
		return "白羊座"
	}
	if d >= 420 && d <= 520 {
		return "金牛座"
	}
	if d >= 521 && d <= 621 {
		return "双子座"
	}
	if d >= 622 && d <= 722 {
		return "巨蟹座"
	}
	if d >= 723 && d <= 822 {
		return "狮子座"
	}
	if d >= 823 && d <= 922 {
		return "处女座"
	}
	if d >= 923 && d <= 1023 {
		return "天秤座"
	}
	if d >= 1024 && d <= 1122 {
		return "天蝎座"
	}

	if d >= 1123 && d <= 1221 {
		return "射手座"
	}
	if d >= 1222 || d <= 119 {
		return "魔羯座"
	}
	if d >= 120 && d <= 218 {
		return "水平座"
	}
	if d >= 219 && d <= 320 {
		return "双鱼座"
	}

	return "水平座"
}