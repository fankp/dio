package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	Context *gin.Context
}

const (
	tokenHeaderName = "Authorization"
)

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (g *Gin) GetToken() string {
	return g.Context.GetHeader(tokenHeaderName)
}

func (g *Gin) GetUserId() int32 {
	userId, exists := g.Context.Get("user_id")
	if exists {
		return userId.(int32)
	}
	return -1
}

func (g *Gin) GetChName() string {
	chName, exists := g.Context.Get("ch_name")
	if exists {
		return chName.(string)
	}
	return ""
}

func (g *Gin) GetOperator() string {
	return fmt.Sprintf("%d,%s", g.GetUserId(), g.GetChName())
}

func (g *Gin) IsAdmin() bool {
	return g.Context.GetBool("admin")
}

func (g *Gin) Response(httpCode int, success bool, msg string, data interface{}) {
	g.Context.JSON(httpCode, &Response{
		Success: success,
		Msg:     msg,
		Data:    data,
	})
}
