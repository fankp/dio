package user

import (
	"context"
	"devops-integral/upm-srv/models/user"
	proto "devops-integral/upm-srv/proto/user"
)

type Handler struct {
}

func (h Handler) GetUserById(ctx context.Context, req *proto.GetUserByIdReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	u, err := userService.QueryByUserId(req.UserId)
	if err != nil {
		// 查询失败，返回异常信息
		resp.Success = false
		resp.Message = err.Error()
		return nil
	}
	// 查询成功，返回用户信息
	resp.Success = true
	resp.User = convert2ProtoUser(*u)
	return err
}

func (h Handler) GetUserByName(ctx context.Context, req *proto.GetUserByNameReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	u, err := userService.QueryByName(req.TenantId, req.Username)
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

func (h Handler) CreateUser(ctx context.Context, req *proto.CreateUserReq, resp *proto.GetUserResp) error {
	userService := user.GetUserService()
	u, err := userService.CreateUser(&user.User{
		TenantId:  req.TenantId,
		Username:  req.Username,
		ChName:    req.ChName,
		Password:  req.Password,
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
		TenantId:  req.TenantId,
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
		UserId:               user.UserId,
		TenantId:             user.TenantId,
		Username:             user.Username,
		ChName:               user.ChName,
		Password:             user.Password,
		Email:                user.Email,
		Phone:                user.Phone,
		Admin:                user.Admin,
		CreatedOn:            user.CreatedOn,
		CreatedBy:            user.CreatedBy,
		UpdatedOn:            user.UpdatedOn,
		UpdatedBy:            user.UpdatedBy,
		DeletedOn:            0,
	}
}