package identity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 测试指定时间之间年龄
func TestAgeAt(t *testing.T) {
	b := time.Date(2016, 1, 1 ,0 ,0 ,0 , 0, time.UTC)
	n := time.Date(2019, 1, 1 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, AgeAt(b, n), 3)

	b = time.Date(2016, 2, 3 ,0 ,0 ,0 , 0, time.UTC)
	n = time.Date(2019, 2, 2 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, AgeAt(b, n), 2)

	b = time.Date(2016, 2, 3 ,0 ,0 ,0 , 0, time.UTC)
	n = time.Date(2019, 2, 3 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, AgeAt(b, n), 3)
}

// 测试年份是不是闰年
func TestIsLeapYear(t *testing.T) {
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	d := time.Date(2016, 1, 1 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, IsLeapYear(d), true)

	d = time.Date(2015, 1, 1 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, IsLeapYear(d), false)

	d = time.Date(1986, 2, 1 ,0 ,0 ,0 , 0, time.UTC)
	assert.Equal(t, IsLeapYear(d), false)
}