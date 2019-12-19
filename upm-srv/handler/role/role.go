package role

import (
	"context"
	"dio/upm-srv/models/role"
	proto "dio/upm-srv/proto/role"
	"github.com/micro/go-micro/util/log"
)

type Handler struct {
}

func (h Handler) CreateRole(ctx context.Context, in *proto.CreateRoleReq, out *proto.RoleResp) error {
	roleService := role.GetRoleService()
	// 根据用户ID查询出当前用户的角色
	roleIds, err := roleService.SelectUserRoleIds(in.CreatorUserId, 0)
	if err != nil {
		log.Errorf("根据用户ID查询用户角色信息失败", err)
	}
	var userRoleId int32
	if len(roleIds) > 0 {
		userRoleId = roleIds[0]
	}
	// 调用role的数据库操作创建角色
	r, err := roleService.CreateRole(&role.Role{
		ParentRoleId: userRoleId,
		RoleType:     in.RoleType,
		RoleName:     in.RoleName,
		RoleDesc:     in.RoleDesc,
		CreatedBy:    in.CreatedBy,
		UpdatedBy:    in.CreatedBy,
	})
	if err != nil {
		// 创建失败，返回异常以及错误信息
		out.Success = false
		out.Message = err.Error()
		return err
	}
	// 创建成功，返回角色信息
	out.Success = true
	out.Role = convert2ProtoRole(r)
	return nil
}
func (h Handler) UpdateRole(ctx context.Context, in *proto.UpdateRoleReq, out *proto.RoleResp) error {
	roleService := role.GetRoleService()
	// 调用role的数据库操作更新角色
	r, err := roleService.UpdateRole(&role.Role{
		RoleId:    in.RoleId,
		RoleName:  in.RoleName,
		RoleDesc:  in.RoleDesc,
		UpdatedBy: in.UpdatedBy,
	})
	if err != nil {
		// 更新失败，返回异常以及错误信息
		out.Success = false
		out.Message = err.Error()
		return err
	}
	// 更新成功，返回角色信息
	out.Success = true
	out.Role = convert2ProtoRole(r)
	return nil
}
func (h Handler) DeleteRole(ctx context.Context, in *proto.DeleteRoleReq, out *proto.DeleteRoleResp) error {
	roleService := role.GetRoleService()
	// 调用role的数据库操作更新角色
	err := roleService.DeleteRole(in.RoleId)
	if err != nil {
		// 删除失败，返回异常以及错误信息
		out.Success = false
		out.Message = err.Error()
		return err
	}
	// todo 删除角色与权限的关联关系，删除角色与用户项目的关联关系
	// 删除成功，返回角色信息
	out.Success = true
	return nil
}
func (h Handler) SelectRolesByName(ctx context.Context, in *proto.SelectRoleReq, out *proto.SelectRolesResp) error {
	roleService := role.GetRoleService()
	roles, err := roleService.SelectRolesByName(in.RoleName)
	if err != nil {
		// 删除失败，查询异常以及错误信息
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.Roles = convert2ProtoRoles(roles)
	return nil
}

func (h Handler) AccessRole(ctx context.Context, in *proto.AccessRoleReq, out *proto.AccessRoleResp) error {
	roleService := role.GetRoleService()
	err := roleService.AccessRole(in.RoleId, in.PrivilegeIds, in.CreatedBy)
	if err != nil {
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	return nil
}

func convert2ProtoRole(role *role.Role) *proto.Role {
	return &proto.Role{
		RoleId:    role.RoleId,
		RoleName:  role.RoleName,
		RoleDesc:  role.RoleDesc,
		CreatedOn: role.CreatedOn,
		CreatedBy: role.CreatedBy,
		UpdatedOn: role.UpdatedOn,
		UpdatedBy: role.UpdatedBy,
		DeletedOn: role.DeletedOn,
	}
}

func convert2ProtoRoles(roles []role.Role) []*proto.Role {
	res := make([]*proto.Role, len(roles))
	for i, each := range roles {
		res[i] = convert2ProtoRole(&each)
	}
	return res
}
