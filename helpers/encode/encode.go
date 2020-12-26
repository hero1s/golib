package encode

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strings"
)

// --------------------
// Encode
// 用gob进行数据编码
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// -------------------
// Decode
// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 用json进行数据编码
//
func EncodeJson(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// -------------------
// Decode
// 用json进行数据解码
//
func DecodeJson(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}

//将一个数据结构转填充另一个数据结构
func ChangeStructByEncodeJson(from interface{}, to interface{}) error {
	data, err := EncodeJson(from)
	if err != nil {
		return err
	}
	return DecodeJson(data, to)
}

func MapStringToStruct(m map[string]interface{}, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, &i)
	if err != nil {
		return err
	}
	return nil
}

func StructToMapString(i interface{}, m map[string]interface{}) error {
	bin, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, &m)
	if err != nil {
		return err
	}
	return nil
}

// 如果有相同的key,会被覆盖
func MergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

// parse params(name=nick&pass=123)
func ParseUrlString(params string) map[string]string {
	paramsMap := map[string]string{}
	for _, param := range strings.Split(params, "&") {
		if !strings.Contains(param, "=") {
			continue
		}
		paramList := strings.Split(param, "=")
		paramsMap[paramList[0]] = paramList[1]
	}
	return paramsMap
}

// 过来结构体空字段，转换json字段的map
func ChangeStructPointToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if !f.IsNil() {
			data[t.Field(i).Tag.Get("json")] = f.Interface()
		}
	}
	return data
}

func ChangeStructToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		data[t.Field(i).Tag.Get("json")] = f.Interface()
	}
	return data
}
