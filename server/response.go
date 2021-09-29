package server

import (
	"github.com/gin-gonic/gin"
)

const (
	OK                  = 0
	BadRequest          = 400
	ErrorRequireLogin   = 403
	InternalServerError = 500
)

var statusText = map[int]string{
	OK:                  "success",
	BadRequest:          "bad request params",
	InternalServerError: "server internal error",
}

func respErr(c *gin.Context, msg string) {
	c.JSON(BadRequest, RespJsonObj{
		Code: BadRequest,
		Msg:  msg,
	})
}

func respJson(c *gin.Context, data interface{}) {
	c.JSON(OK, RespJsonObj{
		Code: OK,
		Msg:  statusText[OK],
		Data: data,
	})
}