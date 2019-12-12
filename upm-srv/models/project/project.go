package project

import (
	"devops-integral/basic/db"
	"fmt"
	"sync"
)

var (
	ser  Service
	lock sync.Mutex
)

type Service interface {
	// 创建项目
	CreateProject(project *Project) (*Project, error)
	// 更新项目
	UpdateProject(project *Project) (*Project, error)
	// 查询项目
	SelectUserProjects(userId int32, projectName string) ([]Project, error)
	// 查询所有项目
	SelectAllProjects() ([]Project, error)
}

type service struct {
}

type Project struct {
	ProjectId    int32 `gorm:"primary_key"`
	ProjectCode  string
	ProjectName  string
	ProjectDesc  string
	ProjectOwner int32
	CreatedOn    int32
	CreatedBy    string
	UpdatedOn    int32
	UpdatedBy    string
	DeletedOn    int32
}

type UserProjectRoleRlat struct {
	UserProjectRoleRaltId int32 `gorm:"primary_key"`
	UserId                int32
	ProjectId             int32
	RoleId                int32
	CreatedOn             int32
	CreatedBy             string
	UpdatedOn             int32
	UpdatedBy             string
	DeletedOn             int32
}

func (s *service) CreateProject(project *Project) (*Project, error) {
	err := db.GetDb().Create(&project).Error
	return project, err
}

func (s *service) UpdateProject(project *Project) (*Project, error) {
	err := db.GetDb().Update(&project).Error
	return project, err
}

func (s *service) SelectUserProjects(userId int32, projectName string) ([]Project, error) {
	var (
		projects []Project
		err      error
	)
	selectByRlatSql := `deleted_on = '0' and project_id in 
			(select project_id from upm_user_project_role_rlat t where t.user_id = ? and t.deleted_on = '0')`
	selectOwnerSql := `deleted_on = '0' and project_owner = ?`
	if projectName != "" {
		projectNameCond := fmt.Sprintf(" and project_name like '%s'", projectName+"%")
		selectByRlatSql += projectNameCond
		selectOwnerSql += projectNameCond
	}
	err = db.GetDb().Where(selectByRlatSql, userId).Or(selectOwnerSql, userId).Find(&projects).Error
	return projects, err
}

func (s *service) SelectAllProjects() ([]Project, error) {
	var (
		projects []Project
		err      error
	)
	err = db.GetDb().Where("delete_on = '0'").Find(&projects).Error
	return projects, err
}

func GetProjectService() Service {
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
