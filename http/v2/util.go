package http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func toMap(v interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	bt, _ := json.Marshal(v)
	json.Unmarshal(bt, &m)
	return m
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	} else if bt, ok := v.([]byte); ok {
		return string(bt)
	} else {
		return fmt.Sprintf("%v", v)
	}
}

func toUrlEncoding(data ...interface{}) string {
	if len(data) == 0 {
		return ""
	}
	urlVals := url.Values{}
	for _, q := range data {
		mp := toMap(q)
		for k, v := range mp {
			val := ""
			switch v.(type) {
			case string:
				val = v.(string)
			case int,int64,[]byte,float32,float64:
				val = toString(v)
			default:
				continue
			}
			urlVals.Add(k, val)
		}
	}
	return urlVals.Encode()
}

func fileExist(p string) bool {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}
	return false
}