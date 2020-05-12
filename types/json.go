package types

import "encoding/json"

// JSONString json 字符串
// 在json序列化的时候，会把json字符串转成 object
type JSONString string

// MarshalJSON 格式化为json字符串的时候，会格式化成 object
func (s JSONString) MarshalJSON() ([]byte, error) {
	var data interface{}
	json.Unmarshal([]byte(s), &data)
	return json.Marshal(data)
}

func (s JSONString) JSON(v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// JSONObj json 对象， 序列化的时候，变成纯字符串
type JSONObj map[string]interface{}

func (m JSONObj) String() JSONString {
	b, _ := json.Marshal(m)
	return JSONString(b)
}

func (m JSONObj) JSON(v interface{}) error {
	return m.String().JSON(v)
}
