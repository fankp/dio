package router

import (
	"dio/basic/common/middleware/cors"
	"dio/basic/common/middleware/jwt"
	"dio/upm-web/handler/privilege"
	"dio/upm-web/handler/project"
	"dio/upm-web/handler/user"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	// 全局使用跨域中间件，否则会报404
	r.Use(cors.Cors())
	// 设置跟路径为upm/v1
	groupV1 := r.Group("/upm/v1")
	// 用户相关的路由
	userGroup := groupV1.Group("/user")
	userGroup.POST("/login", user.Login)
	userGroup.POST("/register", user.Register)
	userGroupNeedLogin := userGroup.Group("")
	userGroupNeedLogin.Use(jwt.JWT())
	userGroupNeedLogin.GET("/privileges/*projectId", privilege.Privileges)
	// 项目相关的路由
	projectGroup := groupV1.Group("/project")
	projectGroup.Use(jwt.JWT())
	projectGroup.POST("/create", project.Create)
	projectGroup.GET("/projects", project.UserProjects)
	return r
}
