package webhook

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {
	a := assert.New(t)

	url := "https://open.feishu.cn/open-apis/bot/v2/hook/cf68991a-0580-4d90-9ab5-2da9ed5f9308"
	msg := "test go send lark log"
	ok, err := SendLark(msg, url)

	a.Equal(ok, true)
	a.Nil(err)
}
