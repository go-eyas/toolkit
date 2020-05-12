package types

import (
	"encoding/json"
	"testing"
)

func TestJSONString(t *testing.T) {
	str := JSONString(`{"demo": true, "num": 123}`)
	data := struct {
		S JSONString
	}{str}
	raw, _ := json.Marshal(data)
	t.Logf("JSONString marshal: %s", string(raw))

	data2 := struct {
		Demo bool
		Num  int
	}{}
	str.JSON(&data2)
	t.Logf("JSONString unmarshal: %#v", data2)
}

func TestJSONObj(t *testing.T) {
	obj := JSONObj{
		"demo": true,
		"num":  123,
	}

	data1 := struct {
		Demo bool
		Num  int
	}{}
	obj.JSON(data1)
	t.Logf("JSONObj to json: %#v", data1)
}
