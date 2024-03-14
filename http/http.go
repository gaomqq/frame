package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int64
	Data interface{}
	Msg  string
}

func Res(c *gin.Context, code int64, data interface{}, msg string) {
	httpCode := http.StatusOK
	if code > 20000 {
		httpCode = http.StatusBadGateway
	}
	c.JSON(httpCode, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
	return
}
