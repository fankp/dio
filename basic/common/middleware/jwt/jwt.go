package jwt

import (
	"dio/basic/common/api"
	"dio/basic/common/message"
	"dio/basic/common/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiG := api.Gin{Context: c}
		token := apiG.GetToken()
		// token为空，直接返回401
		if token == "" {
			apiG.Response(http.StatusUnauthorized, false, message.InvalidTokenError, nil)
			apiG.Context.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				// token过期
				apiG.Response(http.StatusUnauthorized, false, message.InvalidTokenError, nil)
				apiG.Context.Abort()
				return
			default:
				// token校验失败
				apiG.Response(http.StatusUnauthorized, false, message.InvalidTokenError, nil)
				apiG.Context.Abort()
				return
			}
		}
		apiG.Context.Set("username", claims.Username)
		apiG.Context.Set("user_id", claims.UserId)
		apiG.Context.Set("ch_name", claims.ChName)
		apiG.Context.Set("admin", claims.Admin)
		// 校验成功，继续执行
		apiG.Context.Next()
	}
}
