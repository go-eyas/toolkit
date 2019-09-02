package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RData struct {
	Code   int         `json:"-"` // http 状态码
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// resp 回应工具
type resp struct {
	c *gin.Context
}

// R 封装响应数据
func R(c *gin.Context) *resp {
	return &resp{c}
}

const CodeSuccess = 0
const CodeUnknowError = 99999999

// Response 响应数据
func (r resp) Response(data *RData) {
	c := r.c
	c.JSON(data.Code, data)
}

// parse 解析响应数据
func (r resp) Parse(v interface{}) *RData {
	data := &RData{
		Code:   http.StatusOK,
		Msg:    "ok",
		Status: CodeSuccess,
	}
	switch v.(type) {
	case error:
		res := v.(error)
		data.Code = http.StatusInternalServerError
		data.Msg = res.Error()
		data.Status = CodeUnknowError
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
			resStatus = CodeSuccess
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

		data = &RData{
			Code:   resCode.(int),
			Status: resStatus.(int),
			Msg:    resMsg.(string),
			Data:   resData,
		}

	case RData, *RData:
		if b, ok := v.(RData); ok {
			data = &b
		} else {
			data = v.(*RData)
		}
	default:
		data.Data = v
	}

	return data
}

// OK 响应成功
func (r resp) OK(v interface{}) {
	r.Response(&RData{
		Code:   http.StatusOK,
		Msg:    "ok",
		Status: CodeSuccess,
		Data:   v,
	})
}

// Res 通用回应
func (r resp) Res(v interface{}) {
	r.Response(r.Parse(v))
}

// Err 回应错误
func (r resp) Err(v error) {
	r.Res(v)
}

func (r resp) Error(v interface{}) {
	data := r.Parse(v)
	if data.Code == http.StatusOK {
		data.Code = http.StatusInternalServerError
	}
	if data.Msg == "ok" {
		if msg, ok := data.Data.(string); ok {
			data.Msg = msg
			data.Data = gin.H{}
		} else {
			data.Msg = "unknow error"
		}
	}
	if data.Status == CodeSuccess {
		data.Status = CodeUnknowError
	}
	r.Response(data)
}

// Forbidden 回应禁止访问
func (r resp) Forbidden(v error) {
	data := r.Parse(v)
	data.Code = http.StatusForbidden
	r.Response(data)
}
