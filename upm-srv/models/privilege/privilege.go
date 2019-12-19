package privilege

import (
	"dio/basic/db"
	"sync"
)

var (
	ser  Service
	lock sync.Mutex
)

type Service interface {
	// 查询角色拥有的权限清单
	SelectPrivilegeCodes(roleId int32) ([]string, error)
	// 查询所有的权限清单
	SelectAllPrivilegeCodes() ([]string, error)
	// 查询允许的权限群组
	SelectPrivilegeGroups(roleIds []int32) ([]PrivilegeGroup, error)
	// 查询所有的权限群组
	SelectAllPrivilegeGroups() ([]PrivilegeGroup, error)
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

func (s *service) SelectPrivilegeCodes(roleId int32) ([]string, error) {
	var (
		privileges []Privilege
		err        error
	)
	selectSql := "deleted_on='0' and privilege_id in (select privilege_id from upm_role_privilege_rlat t where deleted_on='0' and t.role_id = ?)"
	err = db.GetDb().Table("upm_privilege").Select("privilege_code").Where(selectSql, roleId).Find(&privileges).Error
	privilegeCodes := make([]string, len(privileges))
	for i, privilege := range privileges {
		privilegeCodes[i] = privilege.PrivilegeCode
	}
	return privilegeCodes, err
}

func (s *service) SelectAllPrivilegeCodes() ([]string, error) {
	var (
		privileges []Privilege
	)
	err := db.GetDb().Table("upm_privilege").Select("privilege_code").Where("deleted_on='0'").Find(&privileges).Error
	privilegeCodes := make([]string, len(privileges))
	for i, privilege := range privileges {
		privilegeCodes[i] = privilege.PrivilegeCode
	}
	return privilegeCodes, err
}

func (s *service) SelectPrivilegeGroups(roleIds []int32) ([]PrivilegeGroup, error) {
	var (
		privilegeGroups []PrivilegeGroup
		err             error
	)
	selectPrivilegeSql := "deleted_on='0' and privilege_id in (select privilege_id from upm_role_privilege_rlat t where deleted_on='0' and t.role_id in (?))"
	err = db.GetDb().Preload("Privilege", selectPrivilegeSql, roleIds).Where(&PrivilegeGroup{
		DeletedOn: 0,
	}).Find(&privilegeGroups).Error
	return privilegeGroups, err
}

func (s *service) SelectAllPrivilegeGroups() ([]PrivilegeGroup, error) {
	var (
		privilegeGroups []PrivilegeGroup
		err             error
	)
	err = db.GetDb().Preload("Privilege", &Privilege{
		DeletedOn: 0,
	}).Where(&PrivilegeGroup{
		DeletedOn: 0,
	}).Find(&privilegeGroups).Error
	return privilegeGroups, err
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
