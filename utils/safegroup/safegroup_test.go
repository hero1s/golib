package safegroup

import (
	"testing"

	"github.com/hero1s/golib/internal/utest"
)

func TestGroup_Go(t *testing.T) {
	assert := utest.NewAssert(t, "TestGroup_Go")
	var sg Group
	sg.Go(func() error {
		return nil
	})
	sg.Go(func() error {
		panic("this is a panic")
	})
	err := sg.Wait()
	assert.IsNotNil(err)
}
