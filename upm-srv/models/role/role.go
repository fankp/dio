package role

import (
	"devops-integral/basic/db"
	"sync"
)

var (
	ser  Service
	lock sync.Mutex
)

type Service interface {
	// 创建角色
	CreateRole(role *Role) (*Role, error)
	// 删除角色
	DeleteRole(roleId int32) error
	// 更新角色
	UpdateRole(role *Role) (*Role, error)
	// 查询根据角色名角色
	SelectRolesByName(roleName string) (*[]Role, error)
	// 给角色赋权
	AccessRole(roleId int32, rolePrivilegeRlats []RolePrivilegeRlat) error
}

type service struct {
}

type Role struct {
	RoleId      int32 `gorm:"primary_key"`
	RoleName    string
	RoleDesc    string
	CreatedOn        int32
	CreatedBy        string
	UpdatedOn        int32
	UpdatedBy        string
	DeletedOn        int32
}

type RolePrivilegeRlat struct {
	RolePrivilegeRlatId   int32 `gorm:"primary_key"`
	RoleId int32
	PrivilegeId         int32
	CreatedOn          int32
	CreatedBy          string
	UpdatedOn          int32
	UpdatedBy          string
	DeletedOn          int32
}

func (s service) CreateRole(role *Role) (*Role, error) {
	err := db.GetDb().Create(&role).Error
	return role, err
}

func (s service) DeleteRole(roleId int32) error {
	err := db.GetDb().Delete(&Role{RoleId:roleId}).Error
	return err
}

func (s service) UpdateRole(role *Role) (*Role, error) {
	err := db.GetDb().Update(&role).Error
	return role, err
}

func (s service) SelectRolesByName(roleName string) (*[]Role, error) {
	var (
		roles []Role
		err error
	)
	err = db.GetDb().Where("deleted_on = '0' and role_name like", roleName + "%").Find(&roles).Error
	return &roles, err
}

func (s service) AccessRole(roleId int32, rolePrivilegeRlats []RolePrivilegeRlat) error{
	var err error
	// 开启事物
	tx := db.GetDb().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 判断是否开启事物失败
	if err = tx.Error; err != nil {
		return err
	}
	// 删除数据库中已存在的角色关联关系
	err = tx.Delete(&RolePrivilegeRlat{RoleId:roleId}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(rolePrivilegeRlats) > 0 {
		// 插入到数据库
		for each := range rolePrivilegeRlats {
			err := tx.Create(&each).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func GetPrivilegeService() Service {
	// 设置锁
	lock.Lock()
	defer lock.Unlock()
	if ser == nil {
		// 初始化ser
		ser = &service{}
	}
	// 返回实例
	return ser
}

