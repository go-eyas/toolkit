package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorData struct {
	Code   int         `json:"-"` // http 状态码
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// parse 解析响应数据
func ParseError(v interface{}) *ErrorData {
	data := &ErrorData{
		Code:   http.StatusOK,
		Msg:    "ok",
		Status: 0,
	}
	switch v.(type) {
	case error:
		res := v.(error)
		data.Code = http.StatusInternalServerError
		data.Msg = res.Error()
		data.Status = 999999999
		data.Data = gin.H{}

	case string:
		data.Data = v.(string)

	case gin.H, *gin.H, map[string]interface{}:
		var e gin.H
		if b, ok := v.(gin.H); ok {
			e = b
		} else if b, ok := v.(map[string]interface{}); ok {
			e = gin.H(b)
		} else if b, ok := v.(*gin.H); ok {
			e = *b
		}

		resCode := e["code"]
		if resCode == nil {
			resCode = http.StatusOK
		}

		resStatus := e["status"]
		if resStatus == nil {
			resStatus = 0
		}

		resMsg := e["msg"]
		if resMsg == nil {
			resMsg = "ok"
		} else if errmsgError, ok := resMsg.(error); ok {
			resMsg = errmsgError.Error()
		}

		resData := e["data"]
		if resData == nil {
			resData = gin.H{}
		}

		data = &ErrorData{
			Code:   resCode.(int),
			Status: resStatus.(int),
			Msg:    resMsg.(string),
			Data:   resData,
		}

	case ErrorData, *ErrorData:
		if b, ok := v.(ErrorData); ok {
			data = &b
		} else {
			data = v.(*ErrorData)
		}
	default:
		data.Data = v
	}

	return data
}
