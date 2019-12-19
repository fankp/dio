package basic

import (
	"devops-integral/basic/config"
	"devops-integral/basic/db"
	"devops-integral/basic/redis"
)

func Init() {
	// 加载配置信息
	config.Init()
	//加载数据库连接
	db.Init()
	// 初始化Redis
	redis.Init()
}
