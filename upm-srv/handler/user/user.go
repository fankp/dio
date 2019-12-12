package user

import (
	"context"
	"devops-integral/basic/common/message"
	"devops-integral/basic/common/utils"
	"devops-integral/upm-srv/models/user"
	proto "devops-integral/upm-srv/proto/user"
	"github.com/jinzhu/gorm"
)

// 用户密码加密密钥，必须是16位或24位或32位
const userSecret = "usersecretdevops"

type Handler struct {
}

func (h Handler) GetUserById(ctx context.Context, req *proto.GetUserByIdReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	u, err := userService.QueryByUserId(req.UserId)
	if err != nil {
		// 查询失败，返回异常信息
		resp.Success = false
		resp.Message = err.Error()
		return err
	}
	// 查询成功，返回用户信息
	resp.Success = true
	resp.User = convert2ProtoUser(*u)
	return err
}

func (h Handler) CheckUser(ctx context.Context, req *proto.CheckUserReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	// 根据租户和用户名查询用户
	u, err := userService.QueryByName(req.Username)
	if err != nil {
		// 查询失败，返回异常信息
		resp.Success = false
		if err == gorm.ErrRecordNotFound {
			resp.Message = message.UserNotExists
			return nil
		}
		return err
	}
	// 查询成功，比对用户密码
	if utils.AesEncrypt(req.Password, userSecret) != u.Password {
		// 密码不正确
		resp.Success = false
		resp.Message = message.PasswordError
		return nil
	}
	resp.Success = true
	resp.User = convert2ProtoUser(*u)
	return nil
}

func (h Handler) CreateUser(ctx context.Context, req *proto.CreateUserReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	// 创建用户，对用户密码进行加密
	u, err := userService.CreateUser(&user.User{
		Username:  req.Username,
		ChName:    req.ChName,
		Password:  utils.AesEncrypt(req.Password, userSecret),
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	})
	if err != nil {
		// 查询失败，返回异常信息
		resp.Success = false
		resp.Message = err.Error()
		return err
	}
	// 查询成功，返回用户信息
	resp.Success = true
	resp.User = convert2ProtoUser(*u)
	return nil
}

func (h Handler) UpdateUser(ctx context.Context, req *proto.UpdateUserReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	u, err := userService.UpdateUser(&user.User{
		UserId:    req.UserId,
		Username:  req.Username,
		ChName:    req.ChName,
		Email:     req.Email,
		Phone:     req.Phone,
		UpdatedBy: req.UpdatedBy,
	})
	if err != nil {
		// 查询失败，返回异常信息
		resp.Success = false
		resp.Message = err.Error()
		return err
	}
	// 查询成功，返回用户信息
	resp.Success = true
	// 通过反射把model总的User转换中proto中的User
	resp.User = convert2ProtoUser(*u)
	return nil
}

func convert2ProtoUser(user user.User) *proto.User {
	return &proto.User{
		UserId:    user.UserId,
		Username:  user.Username,
		ChName:    user.ChName,
		Password:  user.Password,
		Email:     user.Email,
		Phone:     user.Phone,
		Admin:     user.Admin,
		CreatedOn: user.CreatedOn,
		CreatedBy: user.CreatedBy,
		UpdatedOn: user.UpdatedOn,
		UpdatedBy: user.UpdatedBy,
		DeletedOn: 0,
	}
}
