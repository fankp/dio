package privilege

import (
	"context"
	"devops-integral/basic/common/api"
	"devops-integral/basic/common/constants"
	"devops-integral/basic/common/message"
	privilege "devops-integral/upm-srv/proto/privilege"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"net/http"
)

var (
	privilegeService = privilege.NewPrivilegeService(constants.ServiceNameUpmSrv, client.DefaultClient)
)

func PRIVILEGE(privilegeCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiG := api.Gin{Context: c}
		// admin用户直接放行
		if apiG.IsAdmin() {
			apiG.Context.Next()
			return
		}
		resp, err := privilegeService.CheckPrivilege(context.TODO(), &privilege.CheckPrivilegeReq{
			UserId:    apiG.GetUserId(),
			ProjectId: 0,
			Admin:     apiG.IsAdmin(),
		})
		if err != nil {
			apiG.Response(http.StatusOK, false, message.CheckPrivilegeError, err.Error())
			apiG.Context.Abort()
			return
		}
		if !resp.Passed {
			apiG.Response(http.StatusOK, false, message.NoPrivilegeError, err.Error())
			apiG.Context.Abort()
			return
		}
		// 校验成功，继续执行
		apiG.Context.Next()
	}
}
