package main

import (
	"dio/basic"
	"dio/basic/common/constants"
	"dio/basic/config"
	"dio/upm-web/router"
	"fmt"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
)

func main() {
	// 初始化运行依赖的相关信息
	log.Logf("开始启动服务，服务名称：%s", constants.ServiceNameUpmWeb)
	basic.Init()
	// 使用etcd进行服务注册
	micReg := etcd.NewRegistry(registryOptions)
	service := web.NewService(
		web.Name(constants.ServiceNameUpmWeb),
		web.Registry(micReg),
		web.Address(":8001"))
	if err := service.Init(); err != nil {
		log.Fatalf("初始化服务失败", err)
	}
	// 使用gin初始化router
	r := router.Init()
	service.Handle("/", r)
	// 启动服务
	if err := service.Run(); err != nil {
		log.Fatalf("启动服务失败：", err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}
