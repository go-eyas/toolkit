package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	data := struct {
		TM Time
	}{Time(time.Now())}

	b, _ := json.Marshal(data)
	t.Logf("test timt format: %s", string(b))
}
