package privilege

import (
	"context"
	"dio/basic/common/api"
	"dio/basic/common/constants"
	"dio/basic/common/message"
	privilege "dio/upm-srv/proto/privilege"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"net/http"
	"strconv"
)

var (
	privilegeService = privilege.NewPrivilegeService(constants.ServiceNameUpmSrv, client.DefaultClient)
)

func Privileges(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	projectId, _ := strconv.ParseInt(apiG.Context.Param("projectId"), 10, 32)
	resp, err := privilegeService.SelectPrivilegeCodes(context.TODO(), &privilege.SelectPrivilegesReq{
		UserId:    apiG.GetUserId(),
		ProjectId: int32(projectId),
		Admin:     apiG.IsAdmin(),
	})
	if err != nil {
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	apiG.Response(http.StatusOK, true, "", resp.PrivilegeCodes)
}
