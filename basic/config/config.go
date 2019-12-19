package config

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
	"path/filepath"
	"strings"
	"sync"
)

var (
	defaultRootPath         = "app"
	defaultConfigFilePrefix = "application-"
	profiles                defaultProfiles
	dbConfig                defaultDbConfig
	etcdConfig              defaultEtcdConfig
	jwtConfig               defaultJwtConfig
	redisConfig             defaultRedisConfig
	lock                    sync.Mutex
	inited                  bool
)

func Init() {
	lock.Lock()
	defer lock.Unlock()
	if inited {
		log.Logf("配置文件已经初始化过，跳过重新加载")
		return
	}
	// 加载配置文件
	// 获取当前工作路径
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("./", string(filepath.Separator))))
	log.Logf("开始加载配置文件，工作目录：%s", appPath)
	configDir := filepath.Join(appPath, "conf")
	// 加载application.yml文件
	if err := config.Load(file.NewSource(file.WithPath(filepath.Join(configDir, "application.yml")))); err != nil {
		log.Errorf("加载配置文件application.yml失败：", err)
		panic(err)
	}
	if err := config.Get(defaultRootPath, "profiles").Scan(&profiles); err != nil {
		log.Errorf("获取profiles失败：", err)
		panic(err)
	}
	// 加载profiles对应的其他文件
	if len(profiles.GetInclude()) > 0 {
		includes := strings.Split(profiles.GetInclude(), ",")
		sources := make([]source.Source, len(includes))
		for i := 0; i < len(includes); i++ {
			filePath := filepath.Join(configDir, defaultConfigFilePrefix+strings.TrimSpace(includes[i])+".yml")
			log.Logf("加载配置文件：%s", filePath)
			sources[i] = file.NewSource(file.WithPath(filePath))
		}
		// 加载profiles对应的文件
		if err := config.Load(sources...); err != nil {
			log.Errorf("加载profiles中的配置文件失败：", err)
			panic(err)
		}
	}
	// 把加载到的配置信息赋值给变量
	_ = config.Get(defaultRootPath, "etcd").Scan(&etcdConfig)
	_ = config.Get(defaultRootPath, "db").Scan(&dbConfig)
	_ = config.Get(defaultRootPath, "jwt").Scan(&jwtConfig)
	_ = config.Get(defaultRootPath, "redis").Scan(&redisConfig)
	// 标记状态为已经初始化完成
	inited = true
}

func GetDbConfig() DbConfig {
	return dbConfig
}
func GetEtcdConfig() EtcdConfig {
	return etcdConfig
}
func GetJwtConfig() JwtConfig {
	return jwtConfig
}

func GetRedisConfig() RedisConfig {
	return redisConfig
}
