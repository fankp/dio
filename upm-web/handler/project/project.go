package project

import (
	"context"
	"dio/basic/common/api"
	"dio/basic/common/constants"
	"dio/basic/common/message"
	project "dio/upm-srv/proto/project"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"net/http"
)

var (
	projectService = project.NewProjectService(constants.ServiceNameUpmSrv, client.DefaultClient)
)

type CreateProjectForm struct {
	ProjectCode string `form:"project_code" binding:"required"`
	ProjectName string `form:"project_name" binding:"required"`
	ProjectDesc string `form:"project_desc"`
}

func Create(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	var createProjectForm CreateProjectForm
	if err := apiG.Context.ShouldBind(&createProjectForm); err != nil {
		// 请求参数不合法
		apiG.Response(http.StatusOK, false, message.InvalidRequestParamError, nil)
		return
	}
	resp, err := projectService.CreateProject(context.TODO(), &project.CreateProjectReq{
		ProjectCode:  createProjectForm.ProjectCode,
		ProjectName:  createProjectForm.ProjectName,
		ProjectDesc:  createProjectForm.ProjectDesc,
		ProjectOwner: apiG.GetUserId(),
		CreatedBy:    apiG.GetOperator(),
	})
	if err != nil {
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	if !resp.Success {
		// 创建失败
		apiG.Response(http.StatusOK, false, message.ServerError, resp.Message)
		return
	}
	// 创建成功，返回创建成功信息
	apiG.Response(http.StatusOK, true, "", resp.Project)
}

type UpdateProjectForm struct {
	ProjectId   int32  `form:"project_id" binding:"required"`
	ProjectCode string `form:"project_code" binding:"required"`
	ProjectName string `form:"project_name" binding:"required"`
	ProjectDesc string `form:"project_desc"`
}

func Update(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	var updateProjectForm UpdateProjectForm
	if err := apiG.Context.ShouldBind(&updateProjectForm); err != nil {
		// 请求参数不合法
		apiG.Response(http.StatusOK, false, message.InvalidRequestParamError, nil)
		return
	}
	resp, err := projectService.UpdateProject(context.TODO(), &project.UpdateProjectReq{
		ProjectId:   updateProjectForm.ProjectId,
		ProjectCode: updateProjectForm.ProjectCode,
		ProjectName: updateProjectForm.ProjectName,
		ProjectDesc: updateProjectForm.ProjectDesc,
		UpdatedBy:   apiG.GetOperator(),
	})
	if err != nil {
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	if !resp.Success {
		// 更新失败
		apiG.Response(http.StatusOK, false, message.ServerError, resp.Message)
		return
	}
	// 更新成功，返回更新成功信息
	apiG.Response(http.StatusOK, true, "", resp.Project)
}

type SelectUserProject struct {
	ProjectName string `form:"project_name"`
}

func UserProjects(ctx *gin.Context) {
	apiG := api.Gin{Context: ctx}
	var selectUserProject SelectUserProject
	if err := apiG.Context.ShouldBind(&selectUserProject); err != nil {
		// 请求参数不合法
		apiG.Response(http.StatusOK, false, message.InvalidRequestParamError, nil)
		return
	}
	resp, err := projectService.SelectUserProjects(context.TODO(), &project.SelectUserProjectsReq{
		UserId:      apiG.GetUserId(),
		ProjectName: selectUserProject.ProjectName,
	})
	if err != nil {
		apiG.Response(http.StatusOK, false, message.ServerError, err.Error())
		return
	}
	if !resp.Success {
		// 更新失败
		apiG.Response(http.StatusOK, false, message.ServerError, resp.Message)
		return
	}
	// 更新成功，返回更新成功信息
	apiG.Response(http.StatusOK, true, "", resp.Projects)
}
