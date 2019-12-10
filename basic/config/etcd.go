package config

type EtcdConfig interface {
	// 获取是否启用
	GetEnable() bool
	// 获取etcd主机
	GetHost() string
	// 获取etcd端口
	GetPort() int
}

type defaultEtcdConfig struct {
	Enable bool   `json:"enable"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func (c defaultEtcdConfig) GetEnable() bool {
	return c.Enable
}

func (c defaultEtcdConfig) GetHost() string {
	return c.Host
}

func (c defaultEtcdConfig) GetPort() int {
	return c.Port
}
