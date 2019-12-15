package privilege

import (
	"devops-integral/basic/db"
	"sync"
)

var (
	ser  Service
	lock sync.Mutex
)

type Service interface {
	// 查询角色拥有的权限清单
	SelectPrivileges(roleId int32) (*[]Privilege, error)
	// 查询所有的权限清单
	SelectAllPrivileges() (*[]Privilege, error)
	// 查询允许的权限群组
	SelectPrivilegeGroups(roleId int32) (*[]PrivilegeGroup, error)
	// 查询所有的权限群组
	SelectAllPrivilegeGroups() (*[]PrivilegeGroup, error)
}

type service struct {
}

type Privilege struct {
	PrivilegeId      int32 `gorm:"primary_key"`
	PrivilegeGroupId int32
	PrivilegeCode    string
	PrivilegeName    string
	CreatedOn        int32
	CreatedBy        string
	UpdatedOn        int32
	UpdatedBy        string
	DeletedOn        int32
}

type PrivilegeGroup struct {
	PrivilegeGroupId   int32 `gorm:"primary_key"`
	PrivilegeGroupName string
	Privileges         []Privilege `gorm:"foreignkey:PrivilegeGroupId;association_foreignkey:PrivilegeGroupId"`
	CreatedOn          int32
	CreatedBy          string
	UpdatedOn          int32
	UpdatedBy          string
	DeletedOn          int32
}

func (s *service) SelectPrivileges(roleId int32) (*[]Privilege, error) {
	var (
		privileges []Privilege
		err error
	)
	selectSql := "deleted_on='0' and privilege_id in (select privilege_id from upm_role_privilege_rlat t where deleted_on='0' and t.role_id = ?)"
	err = db.GetDb().Where(selectSql, roleId).Find(&privileges).Error
	return &privileges, err
}

func (s *service) SelectAllPrivileges() (*[]Privilege, error) {
	var (
		privileges []Privilege
		err error
	)
	err = db.GetDb().Where(&Privilege{DeletedOn:0}).Find(&privileges).Error
	return &privileges, err
}

func (s *service) SelectPrivilegeGroups(roleId int32) (*[]PrivilegeGroup, error) {
	var (
		privilegeGroups []PrivilegeGroup
		err error
	)
	selectPrivilegeSql := "deleted_on='0' and privilege_id in (select privilege_id from upm_role_privilege_rlat t where deleted_on='0' and t.role_id = ?)"
	err = db.GetDb().Preload("Privilege", selectPrivilegeSql, roleId).Where(&PrivilegeGroup{
		DeletedOn: 0,
	}).Find(&privilegeGroups).Error
	return &privilegeGroups, err
}

func (s *service) SelectAllPrivilegeGroups() (*[]PrivilegeGroup, error) {
	var (
		privilegeGroups []PrivilegeGroup
		err error
	)
	err = db.GetDb().Preload("Privilege", &Privilege{
		DeletedOn: 0,
	}).Where(&PrivilegeGroup{
		DeletedOn: 0,
	}).Find(&privilegeGroups).Error
	return &privilegeGroups, err
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
