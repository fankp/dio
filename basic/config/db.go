package config

type DbConfig interface {
	// 获取是否启用
	GetEnable() bool
	// 获取数据库类型
	GetType() string
	// 获取数据库主机IP
	GetHost() string
	// 获取数据库端口
	GetPort() int
	// 获取数据库名称
	GetName() string
	// 获取数据库用户名
	GetUser() string
	// 获取数据库密码
	GetPassword() string
	// 获取数据库最大空闲连接数
	GetMaxIdleConnection() int
	// 获取数据库最大打开的连接数
	GetMaxOpenConnection() int
	//获取表前缀
	GetTablePrefix() string
}

type defaultDbConfig struct {
	Enable            bool   `json:"enable"`
	Type              string `json:"type"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Name              string `json:"name"`
	User              string `json:"user"`
	Password          string `json:"password"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
	TablePrefix       string `json:"tablePrefix"`
}

func (d defaultDbConfig) GetEnable() bool {
	return d.Enable
}

func (d defaultDbConfig) GetType() string {
	return d.Type
}

func (d defaultDbConfig) GetHost() string {
	return d.Host
}

func (d defaultDbConfig) GetPort() int {
	return d.Port
}

func (d defaultDbConfig) GetName() string {
	return d.Name
}

func (d defaultDbConfig) GetUser() string {
	return d.User
}

func (d defaultDbConfig) GetPassword() string {
	return d.Password
}

func (d defaultDbConfig) GetMaxIdleConnection() int {
	return d.MaxIdleConnection
}

func (d defaultDbConfig) GetMaxOpenConnection() int {
	return d.MaxOpenConnection
}

func (d defaultDbConfig) GetTablePrefix() string {
	return d.TablePrefix
}
