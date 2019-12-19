package privilege

import (
	"context"
	"devops-integral/basic/redis"
	"devops-integral/upm-srv/models/privilege"
	"devops-integral/upm-srv/models/role"
	proto "devops-integral/upm-srv/proto/privilege"
	"fmt"
	"github.com/micro/go-micro/util/log"
	"time"
)

const (
	rolePrivilegeCachePrefix = "upm:privileges:"
	rolePrivilegeCacheAdmin  = "admin"
	privilegeCacheExpire     = 60 * time.Minute
)

type Handler struct {
}

func (h Handler) SelectPrivilegeCodes(ctx context.Context, in *proto.SelectPrivilegesReq, out *proto.SelectPrivilegeCodesResp) error {
	var (
		privilegeCodes []string
		err            error
	)
	privilegeService := privilege.GetPrivilegeService()
	if in.Admin {
		privilegeCodes, err = redis.GetRedis().SMembers(rolePrivilegeCachePrefix + rolePrivilegeCacheAdmin).Result()
		if err != nil {
			log.Errorf("从Redis获取权限缓存失败（管理员用户）", err)
		}
		if len(privilegeCodes) == 0 {
			// 管理员，查询所有的权限清单
			privilegeCodes, err = privilegeService.SelectAllPrivilegeCodes()
			if len(privilegeCodes) > 0 {
				_, err = redis.GetRedis().SAdd(rolePrivilegeCachePrefix+rolePrivilegeCacheAdmin, privilegeCodes).Result()
			}
		}
	} else {
		roleService := role.GetRoleService()
		roleIds, _ := roleService.SelectUserRoleIds(in.UserId, in.ProjectId)
		// 把int类型的数组转换为字符串数组
		cacheNames := make([]string, len(roleIds))
		for i, roleId := range roleIds {
			// 获取cache的key
			cacheName := rolePrivilegeCachePrefix + fmt.Sprintf("%d", roleId)
			cacheNames[i] = cacheName
			h.checkLoadPrivilegeCache(roleId)
		}
		if len(cacheNames) > 0 {
			// 对多个角色的权限编码进行求并集
			privilegeCodes, err = redis.GetRedis().SUnion(cacheNames...).Result()
		}
	}
	// 进行查询结果处理
	if err != nil {
		log.Error("查询权限信息失败", err)
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.PrivilegeCodes = privilegeCodes
	return nil
}

func (h Handler) CheckPrivilege(ctx context.Context, in *proto.CheckPrivilegeReq, out *proto.CheckPrivilegeResp) error {
	if in.Admin {
		// 管理员用户，直接返回成功
		out.Success = true
		out.Passed = true
		return nil
	}
	// 查询用户角色信息
	roleService := role.GetRoleService()
	roleIds, err := roleService.SelectUserRoleIds(in.UserId, in.ProjectId)
	if err != nil {
		log.Errorf("查询用户角色信息失败，用户ID：%d, 项目ID：%，", in.UserId, in.ProjectId, err)
		// 管理员用户，直接返回成功
		out.Success = false
		out.Message = err.Error()
		out.Passed = false
		return err
	}
	for _, roleId := range roleIds {
		h.checkLoadPrivilegeCache(roleId)
		cacheName := rolePrivilegeCachePrefix + fmt.Sprintf("%d", roleId)
		if redis.GetRedis().SIsMember(cacheName, in.PrivilegeCode).Val() {
			out.Success = true
			out.Passed = true
			return nil
		}
	}
	out.Success = true
	out.Passed = false
	return nil
}

func (h Handler) checkLoadPrivilegeCache(roleId int32) {
	privilegeService := privilege.GetPrivilegeService()
	cacheName := rolePrivilegeCachePrefix + fmt.Sprintf("%d", roleId)
	num, err := redis.GetRedis().SCard(cacheName).Result()
	if err != nil {
		log.Errorf("根据Key获取权限缓存元素数量失败，", err)
		return
	}
	if num == 0 {
		privilegeCodes, err := privilegeService.SelectPrivilegeCodes(roleId)
		if err != nil {
			log.Errorf("根据角色查询权限编码失败，", err)
			return
		}
		_, err = redis.GetRedis().SAdd(cacheName, privilegeCodes, privilegeCacheExpire).Result()
		if err != nil {
			log.Errorf("保存权限编码到Redis失败，", err)
			return
		}
	}
}

func (h Handler) SelectPrivilegeGroups(ctx context.Context, in *proto.SelectPrivilegesReq, out *proto.SelectPrivilegeGroupsResp) error {
	var (
		privilegeGroups []privilege.PrivilegeGroup
		err             error
	)
	privilegeService := privilege.GetPrivilegeService()
	if in.Admin {
		// 管理员，查询所有的权限清单
		privilegeGroups, err = privilegeService.SelectAllPrivilegeGroups()
	} else {
		// 普通用户，根据角色ID进行查询
		roleService := role.GetRoleService()
		roleIds, _ := roleService.SelectUserRoleIds(in.UserId, in.ProjectId)
		privilegeGroups, err = privilegeService.SelectPrivilegeGroups(roleIds)
	}
	// 进行查询结果处理
	if err != nil {
		log.Error("查询权限群组信息失败", err)
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.PrivilegeGroups = covert2ProtoPrivilegeGroups(privilegeGroups)
	return nil
}

func covert2ProtoPrivilege(pri *privilege.Privilege) *proto.Privilege {
	return &proto.Privilege{
		PrivilegeId:      pri.PrivilegeId,
		PrivilegeGroupId: pri.PrivilegeGroupId,
		PrivilegeCode:    pri.PrivilegeCode,
		PrivilegeName:    pri.PrivilegeName,
		CreatedOn:        pri.CreatedOn,
		CreatedBy:        pri.CreatedBy,
		UpdatedOn:        pri.UpdatedOn,
		UpdatedBy:        pri.UpdatedBy,
		DeletedOn:        pri.DeletedOn,
	}
}

func covert2ProtoPrivileges(pris []privilege.Privilege) []*proto.Privilege {
	res := make([]*proto.Privilege, len(pris))
	for i, each := range pris {
		res[i] = covert2ProtoPrivilege(&each)
	}
	return res
}

func covert2ProtoPrivilegeGroup(prig *privilege.PrivilegeGroup) *proto.PrivilegeGroup {
	return &proto.PrivilegeGroup{
		PrivilegeGroupId:   prig.PrivilegeGroupId,
		PrivilegeGroupName: prig.PrivilegeGroupName,
		Privileges:         covert2ProtoPrivileges(prig.Privileges),
		CreatedOn:          prig.CreatedOn,
		CreatedBy:          prig.CreatedBy,
		UpdatedOn:          prig.UpdatedOn,
		UpdatedBy:          prig.UpdatedBy,
		DeletedOn:          prig.DeletedOn,
	}
}

func covert2ProtoPrivilegeGroups(prig []privilege.PrivilegeGroup) []*proto.PrivilegeGroup {
	res := make([]*proto.PrivilegeGroup, len(prig))
	for i, each := range prig {
		res[i] = covert2ProtoPrivilegeGroup(&each)
	}
	return res
}
