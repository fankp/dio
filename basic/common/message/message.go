package message

const (
	InvalidRequestParamError = "请求参数不合法"
	ServerError              = "服务异常"
	UserNotExists            = "用户不存在"
	PasswordError            = "密码不正确"
	CreateTokenError         = "创建Token失败"
	InvalidTokenError        = "Token无效或不存在"
	CheckPrivilegeError      = "校验权限信息失败"
	NoPrivilegeError         = "权限不足，请联系管理员"
)
