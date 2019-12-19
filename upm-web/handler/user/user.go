package user

import (
	"context"
	"dio/basic/common/api"
	"dio/basic/common/constants"
	"dio/basic/common/message"
	"dio/basic/common/utils"
	user "dio/upm-srv/proto/user"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/util/log"
	"net/http"
)

var (
	userService = user.NewUserService(constants.ServiceNameUpmSrv, client.DefaultClient)
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	var loginForm LoginForm
	if err := apiG.Context.ShouldBind(&loginForm); err != nil {
		apiG.Response(http.StatusOK, false, message.InvalidRequestParamError, err.Error())
		return
	}
	// 调用rpc服务根据用户名查询用户信息
	respU, err := userService.CheckUser(context.TODO(), &user.CheckUserReq{
		Username: loginForm.Username,
		Password: loginForm.Password,
	})
	if err != nil {
		log.Errorf("调用校验密码服务失败，", err)
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	// 校验密码是否正确
	if !respU.Success {
		log.Errorf("登录失败，用户名：%s，错误信息：%s", loginForm.Username, respU.Message)
		apiG.Response(http.StatusOK, false, respU.Message, nil)
		return
	}
	// 生成token
	token, err := utils.GenerateToken(respU.User.UserId, respU.User.Username, respU.User.ChName, respU.User.Admin)
	if err != nil {
		// 创建token失败
		log.Errorf("生成Token失败", err)
		apiG.Response(http.StatusOK, false, message.CreateTokenError, err.Error())
		return
	}
	apiG.Response(http.StatusOK, true, "", token)
}

type CreateUserForm struct {
	Username string `form:"username" binding:"required"`
	ChName   string `form:"ch_name" binding:"required"`
	Password string `form:"password" binding:"required"`
	Email    string `from:"email"`
	Phone    string `form:"phone"`
}

func Register(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	var createUserFrom CreateUserForm
	if err := apiG.Context.ShouldBind(&createUserFrom); err != nil {
		apiG.Response(http.StatusOK, false, message.InvalidRequestParamError, err.Error())
		return
	}
	resp, err := userService.CreateUser(ctx, &user.CreateUserReq{
		Username:  createUserFrom.Username,
		ChName:    createUserFrom.ChName,
		Password:  createUserFrom.Password,
		Email:     createUserFrom.Email,
		Phone:     createUserFrom.Phone,
		CreatedBy: apiG.GetOperator(),
		UpdatedBy: apiG.GetOperator(),
	})
	if err != nil {
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	apiG.Response(http.StatusOK, true, "", resp.User)
}
