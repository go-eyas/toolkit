package util_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/go-eyas/toolkit/util"
)

func TestAesEecode(t *testing.T) {
	password := []byte("asdfghjkqwertyui")
	var min int64 = 0
	var max int64 = 100
	for i := min; i < max+min; i++ {
		result, err := util.AesEncrypt([]byte(fmt.Sprintf("%d", i)), password)
		if err != nil {
			t.Fatal(err)
		}
		// t.Log(hex.EncodeToString(result)) // 16进制
		t.Log(base64.StdEncoding.EncodeToString(result)) // 16进制
	}

}
