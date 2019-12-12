package user

import (
	"devops-integral/basic/db"
	"sync"
)

var (
	ser  Service
	lock sync.Mutex
)

type Service interface {
	// 根据用户ID查询用户信息
	QueryByUserId(userId int32) (*User, error)
	// 根据用户名称查询用户
	QueryByName(username string) (*User, error)
	// 新增用户
	CreateUser(user *User) (*User, error)
	// 更新用户
	UpdateUser(user *User) (*User, error)
}

type service struct {
}

type User struct {
	UserId    int32 `gorm:"primary_key"`
	Username  string
	ChName    string
	Password  string
	Email     string
	Phone     string
	Admin     bool
	CreatedOn int32
	CreatedBy string
	UpdatedOn int32
	UpdatedBy string
	DeletedOn int32
}

func (s service) QueryByUserId(userId int32) (*User, error) {
	var user = &User{}
	err := db.GetDb().Where(&User{
		UserId:    userId,
		DeletedOn: 0,
	}).First(user).Error
	return user, err
}

func (s service) QueryByName(username string) (*User, error) {
	var user = &User{}
	err := db.GetDb().Where(User{
		Username:  username,
		DeletedOn: 0,
	}).First(user).Error
	return user, err
}

func (s service) CreateUser(user *User) (*User, error) {
	err := db.GetDb().Create(&user).Error
	return user, err
}

func (s service) UpdateUser(user *User) (*User, error) {
	err := db.GetDb().Where("user_id = ?", user.UserId).Update(&user).Error
	return user, err
}

func GetUserService() Service {
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
