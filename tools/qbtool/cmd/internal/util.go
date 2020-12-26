package internal

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	SwaggerVersion = "3"
	SwaggerLink    = "https://github.com/beego/swagger/archive/v" + SwaggerVersion + ".zip"
)

func DownloadFromURL(url, filename string) error {
	var down bool
	if fd, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		down = true
	} else if fd.Size() == int64(0) {
		down = true
	} else {
		log.Printf("%s already exists\n", filename)
		return errors.New("already exist")
	}
	if down {
		log.Printf("Downloading %s to %s...\n", url, filename)
		output, err := os.Create(filename)
		if err != nil {
			log.Printf("Error while creating %s: %s\n", filename, err)
			return err
		}
		defer output.Close()

		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error while downloading %s: %s\n", url, err)
			return err
		}
		defer response.Body.Close()

		n, err := io.Copy(output, response.Body)
		if err != nil {
			log.Printf("Error while downloading %s: %s\n", url, err)
			return err
		}
		log.Printf("%d bytes downloaded!\n", n)
	}
	return nil
}

func UnzipAndDelete(realDst, src string) error {
	log.Printf("Unzipping %s...\n", src)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	rp := strings.NewReplacer("swagger-"+SwaggerVersion, "swagger")
	var nowDir string
	for _, f := range r.File {
		if err := doFile(rp, f); err != nil {
			return err
		}
		if f.FileInfo().IsDir() {
			nowDir = rp.Replace(f.Name)
		}
	}
	if realDst != "." && realDst != strings.TrimSuffix(nowDir, "/") {
		os.RemoveAll(realDst)
		os.Rename(nowDir, realDst)
	}
	log.Printf("Done! Deleting %s...\n", src)
	return os.RemoveAll(src)
}

func doFile(rp *strings.Replacer, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	fName := rp.Replace(f.Name)
	if f.FileInfo().IsDir() {
		os.MkdirAll(fName, f.Mode())
	} else {
		f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
