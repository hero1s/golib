package alisdk

import (
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/utils/strutil"
	"io"
	"strings"
	"time"
)

func UploadFileToOss(reader io.Reader, filename, savePath, savename string, id uint64) (ossPath string, err error) {
	ext := strutil.SubString(strutil.Unicode(filename),
		strings.LastIndex(strutil.Unicode(filename), "."), 5)
	filename = fmt.Sprintf("%d%d", time.Now().UnixNano(), id) + ext
	if len(savename) > 1 {
		filename = savename + ext
	}
	ossPath = GetImagePath(savePath, filename)
	ossPath = strings.ReplaceAll(ossPath, "\\", "/")
	err = PutFileStreamToOss(ossPath, reader)
	log.Debugf("upload file to oss:%v/%v", savePath, filename)
	return ossPath, err

}
