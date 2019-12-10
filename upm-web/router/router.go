package router

import (
	"devops-integral/upm-web/handler/user"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.POST("/upm/login", user.Login)
	r.POST("/upm/user/create", user.CreateUser)
	return r
}
