package config

type JwtConfig interface {
	GetEnable() bool
	GetSecret() []byte
}

type defaultJwtConfig struct {
	Enable bool   `json:"enable"`
	Secret []byte `json:"secret"`
}

func (j defaultJwtConfig) GetEnable() bool {
	return j.Enable
}

func (j defaultJwtConfig) GetSecret() []byte {
	return j.Secret
}
