
package cache

import (
	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = bm.Put("astaxie", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}
	if v := bm.Get("astaxie"); v.(string) != "author" {
		t.Error("get err")
	}

	if err = bm.Put("astaxie1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}

	//test string
	if err = bm.Put("astaxie", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}
	if v := bm.Get("astaxie"); v.(string) != "author" {
		t.Error("get err")
	}

	//test GetMulti
	if err = bm.Put("astaxie1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	os.RemoveAll("cache")
}
