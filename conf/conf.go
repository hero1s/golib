package conf

import (
	"fmt"
	"github.com/hero1s/golib/helpers/file"
	"github.com/hero1s/golib/log"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	json = iota + 1
	yaml
	toml
)

func AutoParseFile(fileName string, c interface{}) error {
	ext := file.GetExt(fileName)
	switch ext {
	case ".json":
		return ParseJson(fileName, c)
	case ".toml":
		return ParseToml(fileName, c)
	case ".yaml":
		return ParseYaml(fileName, c)
	case ".yml":
		return ParseYaml(fileName, c)
	}
	log.Error("文件类型错误:%s", ext)
	return errors.New(fmt.Sprintf("文件类型错误:%s", ext))
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseYaml(confPath string, c interface{}) error {
	return parse(confPath, yaml, c)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseToml(confPath string, c interface{}) error {
	return parse(confPath, toml, c)
}

// 解析配置文件
//    c: 需要解析的相对应的结构体指针，例：conf_test.go
func ParseJson(confPath string, c interface{}) error {
	return parse(confPath, json, c)
}

func parse(confPath string, cType int, c interface{}) error {
	if c == nil {
		return errors.New("c struct ptr can not be nil")
	}
	var cFileType string
	beanValue := reflect.ValueOf(c)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("c must be ptr")
	}
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("c must be struct ptr")
	}
	if confPath == "" {
		return errors.New("load config file path failed, add arguments -conf ")
	}
	viper.SetConfigFile(confPath)
	switch cType {
	case json:
		cFileType = "json"
	case yaml:
		cFileType = "yaml"
	case toml:
		cFileType = "toml"
	default:
		return errors.New("config file only support: yaml、json、toml")
	}
	viper.SetConfigType(cFileType)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}
