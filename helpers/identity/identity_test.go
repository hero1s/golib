package identity

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var tests = map[string]string {
	"18": "511521198602035558",
	"15": "511521860203555",
}

// 测试解析15位身份证
func TestNew15(t *testing.T) {
	s := tests["15"]
	id := New(s)
	err := id.Parse()
	assert.Equal(t, err, nil)
	assert.Equal(t, id.GetProvince(), "51")
	assert.Equal(t, id.GetCity(), "15")
	assert.Equal(t, id. GetArea(), "21")

	assert.Equal(t, id.GetBirthday(), "19860203")
	assert.Equal(t, id.GetSex(), SexMale)

	assert.Equal(t, id.GetBirthdayTime().Year(), 1986)
	assert.Equal(t, id.GetBirthdayTime().Month(), time.Month(2))
	assert.Equal(t, id.GetBirthdayTime().Day(), 3)
	assert.Equal(t, AgeAt(id.GetBirthdayTime(), time.Date(2019, 3, 6, 0, 0, 0, 0, time.UTC)), 33)

	assert.Equal(t, id.GetSequenceCode(), "555")
	assert.Equal(t, id.GetVerifyCode(), "")
}

func TestNew18(t *testing.T) {
	s := tests["18"]
	id := New(s)
	err := id.Parse()
	assert.Equal(t, err, nil)
	assert.Equal(t, id.GetProvince(), "51")
	assert.Equal(t, id.GetCity(), "15")
	assert.Equal(t, id. GetArea(), "21")

	assert.Equal(t, id.GetBirthday(), "19860203")
	assert.Equal(t, id.GetSex(), SexMale)
	assert.Equal(t, id.GetSequenceCode(), "555")
	assert.Equal(t, id.GetVerifyCode(), "8")

}
