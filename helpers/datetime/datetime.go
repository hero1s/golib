package datetime

import (
	"errors"
	"time"
)

var (
	/** 设置每周的起始时间 */
	WeekStartDay  = time.Sunday

	/** 指定日期和时间的默认转换格式 */
	DateTimeFormats = []string{"1/2/2006","1/2/2006 15:4:5","2006","2006-1-2","2006-01-02 15:04:05","20060102150405","15:4:5 Jan 2, 2006 MST"}
)

const (
	DefalutFormat  = "2006-01-02 15:04:05"
	CompressFormat = "20060102150405"
)

// DateTime 结构体
type DateTime struct {
	time.Time
}

// 当前秒
func CurrentSecond() time.Time {
	return time.Now().Truncate(time.Second)
}

// 当前分钟
func CurrentMinute() time.Time {
	return time.Now().Truncate(time.Minute)
}

// 当前小时
func CurrentHour() time.Time {
	return time.Now().Truncate(time.Hour)
}

// 返回今天的日期
func Today() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
}

// 返回今天的最后一刻
func TodayEndMoment() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), time.Now().Location())
}

// 返回本周的第一刻
func BeginThisWeek() time.Time {
	today := Today()
	weekday := int(today.Weekday())
	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		}else {
			weekday = weekday - weekStartDayInt
		}
	}
	return today.AddDate(0,0,-weekday)
}

// 返回本周最后一刻
func EndThisWeek() time.Time {
	return BeginThisWeek().AddDate(0,0,7).Add(-time.Nanosecond)
}

// 返回本月的第一刻
func BeginThisMonth() time.Time {
	year, month, _ := time.Now().Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
}

// 返回本月最后一刻
func EndThisMonth() time.Time {
	return BeginThisMonth().AddDate(0,1,0).Add(-time.Nanosecond)
}

// 返回本年的第一刻
func BeginThisYear() time.Time {
	year, _, _ := time.Now().Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
}

// 返回本年的最后一刻
func EndThisYear() time.Time {
	return BeginThisYear().AddDate(1,0,0).Add(-time.Nanosecond)
}

// 字符串转时间
func Parse(str string) (t time.Time, err error) {
	for _, format := range DateTimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return t,err
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return t,err
}

// 获取当前时间字符串 - yyyy-MM-dd HH:mm:ss
func CurrentDeflult() string {
	return time.Now().Format(DefalutFormat)
}

// 获取当前时间字符串 - yyyyMMddHHmmss
func CurrentCompress() string {
	return time.Now().Format(CompressFormat)
}

//获取当天指定时间点时间戳
func GetCurDayTime(hour, min, second int64) int64 {
	curtime := time.Now().Unix()
	tm := time.Unix(curtime, 0)
	daytime := tm.Hour()*3600 + tm.Minute()*60 + tm.Second()
	addtime := hour*3600 + min*60 + second
	return curtime + int64(addtime) - int64(daytime)
}

//获取给出时间戳的当天的开始与结束的时间戳
func StartEndTimeByTimestamp(t int64) (int64, int64) {
	tt := time.Unix(t, 0)
	ts := time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, tt.Location())
	te := ts.AddDate(0, 0, 1)
	return ts.Unix(), te.Unix()
}
