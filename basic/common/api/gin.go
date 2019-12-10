package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Gin struct {
	Context *gin.Context
}

const (
	tenantHeaderName = "Tenant-Id"
	tokenHeaderName  = "Authorization"
)

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (g *Gin) GetTenantId() int32 {
	// 从header中获取租户信息
	tenant := g.Context.GetHeader(tenantHeaderName)
	if tenant != "" {
		// 转换为int类型
		tenantId, _ := strconv.ParseInt(tenant, 10, 32)
		return int32(tenantId)
	}
	// 租户信息不存在，返回默认租户0
	return 0
}

func (g *Gin) Response(httpCode int, success bool, msg string, data interface{}) {
	g.Context.JSON(httpCode, &Response{
		Success: success,
		Msg:     msg,
		Data:    data,
	})
}
