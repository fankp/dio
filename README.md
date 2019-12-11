# 工程依赖
* [https://go-micro.dev/]
* [https://gorm.io/]
# 启动服务
* 启动upm-srv：`cd upm-srv && go run main.go plugin.go `
* 启动upm-web：`cd ../upm-web && go run main.go plugin.go `
* 启动网关：`./micro --registry etcd --registry_address=localhost:2379 --api_namespace devops.integral.web api --handler web`
