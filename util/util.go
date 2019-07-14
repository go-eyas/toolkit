package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"runtime"
)

// Assert 断言 err != nil
func Assert(err error, msg interface{}) {
	if err != nil {
		panic(msg)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$@")

// RandomStr 生成随机字符串
func RandomStr(length int) string {
	var lenthLetter = len(letterRunes)

	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(lenthLetter)]
	}
	return string(b)
}

// Base64Encoding base64 编码
func Base64Encoding(str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded
}

// Base64Decoding base64 解码
func Base64Decoding(enc string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// FuncName 获取函数的名字
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// AssignMap 合并多个map
func AssignMap(maps ...map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	for _, mp := range maps {
		for key, val := range mp {
			m[key] = val
		}
	}

	return m
}

// ToString 把能转成字符串的都转成JSON字符串
func ToString(v interface{}) string {
	bt, _ := json.Marshal(v)
	return string(bt)
}

// HasFile 是否存在该文件
func HasFile(f string) bool {
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return true
	}
	return false
}

// StructToMap 把结构体转成map，key使用json定义的key
func StructToMap(v interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	bt, _ := json.Marshal(v)
	_ = json.Unmarshal(bt, &data)
	return data
}

// ToStruct 把一个结构体转成另一个结构体，以json key作为关联
func ToStruct(raw interface{}, v interface{}) error {
	var err error
	var bt []byte
	if sraw, ok := raw.(string); ok {
		bt = []byte(sraw)
	} else if braw, ok := raw.([]byte); ok {
		bt = braw
	} else {
		bt, err = json.Marshal(raw)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(bt, v)
}

// ByteToReader 将字节转换成读取流
func ByteToReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}

// ByteToReadCloser 将字节转换成一次性的读取流
func ByteToReadCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(ByteToReader(b))
}
