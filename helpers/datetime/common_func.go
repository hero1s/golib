package datetime

import "time"

//是否同一天
func IsSameDay(oldTime, newTime int64) bool {
	tm1 := time.Unix(oldTime, 0)
	tm2 := time.Unix(newTime, 0)
	if tm1.YearDay() == tm2.YearDay() && tm1.Year() == tm2.Year() {
		return true
	}
	return false
}

//判断当前时间是否为周一至周五
func IsWorkday() bool {
	if time.Monday <= time.Now().Weekday() && time.Friday >= time.Now().Weekday() {
		return true
	}
	return false
}

//默认获取今天的零时时间戳
//-1表示前一天的零时时间戳
//1表示第二天的零时时间戳
func SpecifyDayZeroTimestamp(n ...int) int64 {
	t := time.Now()
	tm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	if len(n) >= 1 {
		tm1 := tm.AddDate(0, 0, n[0])
		return tm1.Unix()
	}
	return tm.Unix()
}

//获取现在距离day天前(后)的时间戳
func GetSpecifyDayByTimestamp(t int64, day int) int64 {
	tt := time.Unix(t, 0)
	tm := tt.AddDate(0, 0, day)
	return tm.Unix()
}

//根据时间戳获取当月的第一天
func GetFirstDayByTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	t1 := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
	return t1.Format("20060102")
}

//判断给出的时间戳是否是当月的第一天
func IsFirstDayOnMonth(timestamp int64) bool {
	//统计数据都是第二天统计的，所以要减去一天
	t := time.Unix(timestamp, 0).AddDate(0, 0, -1)
	firstDay := time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
	firstDayStart := firstDay.Unix()
	firstDayEnd := firstDay.AddDate(0, 0, 1).Unix()
	if t.Unix() >= firstDayStart && t.Unix() < firstDayEnd {
		return true
	}
	return false
}

/*
根据时间戳获取日期
input:2018-01-12 16:24:01 这个日期的时间戳:1515745441
output: 20180112
*/
//根据时间戳获取格式化的日期
func GetDayByTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("20060102")
}

//相对于给出的时间戳获取昨天的时间
func GetYesDayByTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).AddDate(0, 0, -1).Format("20060102")
}

/*
根据时间戳获取日期
input:2018-01-12 16:24:01 这个日期的时间戳:1515745441
output: 20180112
*/
//根据时间戳获取格式华的日期
func GetMonthByTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("200601")
}

//判断timestamp是否为当天的时间戳
func IsToday(timestamp int64) bool {
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	t2 := t1.AddDate(0, 0, 1)
	if timestamp >= t1.Unix() && timestamp < t2.Unix() {
		return true
	}
	return false
}

//判断timestamp是否当月的时间戳
func IsCurrentMonth(timestamp int64) bool {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)
	if timestamp >= start.Unix() && timestamp < end.Unix() {
		return true
	}
	return false
}

//判断timestamp是否在相同的月份的时间戳
func IsSameMonth(t1, t2 int64) bool {
	return time.Unix(t1, 0).Format("200601") == time.Unix(t2, 0).Format("200601")
}

//判断timestamp是否上一个月的时间戳
func IsLastMonth(timestamp int64) bool {
	t := time.Now()
	this := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	last := this.AddDate(0, -1, 0)
	if timestamp < this.Unix() && timestamp >= last.Unix() {
		return true
	}
	return false
}

//获取现在离今天结束的时间还有多久
func GetTimeFromTodayEnd() int64 {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := today.AddDate(0, 0, 1)
	return tomorrow.Unix() - t.Unix()
}

//获取现在离本月结束的时间还有多久
func GetTimeFromMonthEnd() int64 {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)
	return end.Unix() - t.Unix()
}

//获取当前月份,输出格式为:1801表示，2018年1月份
func CurrentYearMonth() string {
	return time.Now().Format("0601")
}

//昨天
func Yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("20060102")
}

//两天前
func TwoDayBefore() string {
	return time.Now().AddDate(0, 0, -2).Format("20060102")
}

//获取本月(减去一天,然后计算出昨天的月份)(如果不减去一天，在月份交际之际会计算出错误的月份)
func ThisMonthByTimestamp(t int64) string {
	return time.Unix(t, 0).Format("200601")
}

//上个月(数据是今天凌晨统计昨天的,并把数据加到昨天所在到月分里,所以要减去一天)
func LastMonthByTimestamp(t int64) string {
	return time.Unix(t, 0).AddDate(0, -1, 0).Format("200601")
}

//获取本月(减去一天,然后计算出昨天的月份)(如果不减去一天，在月份交际之际会计算出错误的月份)
func ThisMonth() string {
	return time.Now().AddDate(0, 0, -1).Format("200601")
}

//上个月(数据是今天凌晨统计昨天的,并把数据加到昨天所在到月分里,所以要减去一天)
func LastMonth() string {
	return time.Now().AddDate(0, 0, -1).AddDate(0, -1, 0).Format("200601")
}

//两个月前(数据是今天凌晨统计昨天的,并把数据加到昨天所在到月分里,所以要减去一天)
func TwoMonthBefore() string {
	return time.Now().AddDate(0, 0, -1).AddDate(0, -2, 0).Format("200601")
}

//获取昨天开始与结束的时间戳
func GetYesterdayTimestamp() (int64, int64) {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesterday := today.AddDate(0, 0, -1)
	return yesterday.Unix(), today.Unix()
}

//获取今天开始与结束的时间戳
func GetTodayTimestamp() (int64, int64) {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tomorrow := today.AddDate(0, 0, 1)
	return today.Unix(), tomorrow.Unix()
}

//获取给出当周的时间戳，当周的开始与结束的时间戳
func StartEndTimeByWeek(timestamp int64) (int64, int64) {
	t := time.Unix(timestamp, 0)
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	end := start.AddDate(0, 0, 7)
	return start.Unix(), end.Unix()
}

//获取给出当月的时间戳，当月的开始与结束的时间戳
func StartEndTimeByMonth(timestamp int64) (int64, int64) {
	t := time.Unix(timestamp, 0)
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)
	return start.Unix(), end.Unix()
}
