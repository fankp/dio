package router

import (
	"devops-integral/basic/common/middleware/jwt"
	"devops-integral/upm-web/handler/project"
	"devops-integral/upm-web/handler/user"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	// 设置跟路径为upm/v1
	groupV1 := r.Group("/upm/v1")
	// 用户相关的路由
	userGroup := groupV1.Group("/user")
	userGroup.POST("/login", user.Login)
	userGroup.POST("/register", user.Register)
	// 项目相关的路由
	projectGroup := groupV1.Group("/project")
	projectGroup.Use(jwt.JWT())
	projectGroup.POST("/create", project.Create)
	projectGroup.GET("/projects", project.UserProjects)
	return r
}
