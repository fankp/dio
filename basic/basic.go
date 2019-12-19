package basic

import (
	"dio/basic/config"
	"dio/basic/db"
	"dio/basic/redis"
)

func Init() {
	// 加载配置信息
	config.Init()
	//加载数据库连接
	db.Init()
	// 初始化Redis
	redis.Init()
}
