package types

import (
	"time"

	"encoding/json"
)

// Time 时间别名，在json序列化的时候，会格式成 2006-01-02 15:04:05 这种时间格式
type Time time.Time

func (tm Time) MarshalJSON() ([]byte, error) {
	s := time.Time(tm).Format("2006-01-02 15:04:05")
	return json.Marshal(s)
}
