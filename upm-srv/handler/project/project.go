package project

import (
	"context"
	"dio/upm-srv/models/project"
	proto "dio/upm-srv/proto/project"
)

type Handler struct {
}

// 创建项目
func (h Handler) CreateProject(ctx context.Context, in *proto.CreateProjectReq, out *proto.ProjectResp) error {
	projectService := project.GetProjectService()
	pro, err := projectService.CreateProject(&project.Project{
		ProjectCode:  in.ProjectCode,
		ProjectName:  in.ProjectName,
		ProjectDesc:  in.ProjectDesc,
		ProjectOwner: in.ProjectOwner,
		CreatedBy:    in.CreatedBy,
		UpdatedBy:    in.CreatedBy,
	})
	if err != nil {
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.Project = convert2ProtoProject(pro)
	return nil
}

// 创建项目
func (h Handler) UpdateProject(ctx context.Context, in *proto.UpdateProjectReq, out *proto.ProjectResp) error {
	projectService := project.GetProjectService()
	pro, err := projectService.UpdateProject(&project.Project{
		ProjectId:   in.ProjectId,
		ProjectCode: in.ProjectCode,
		ProjectName: in.ProjectName,
		ProjectDesc: in.ProjectDesc,
	})
	if err != nil {
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.Project = convert2ProtoProject(pro)
	return nil
}

func (h Handler) SelectUserProjects(ctx context.Context, in *proto.SelectUserProjectsReq, out *proto.UserProjectsResp) error {
	projectService := project.GetProjectService()
	projects, err := projectService.SelectUserProjects(in.UserId, in.ProjectName)
	if err != nil {
		out.Success = false
		out.Message = err.Error()
		return err
	}
	out.Success = true
	out.Projects = convert2ProtoProjects(projects)
	return nil
}

func convert2ProtoProjects(projects []project.Project) []*proto.Project {
	res := make([]*proto.Project, len(projects))
	for i, each := range projects {
		res[i] = convert2ProtoProject(&each)
	}
	return res
}

func convert2ProtoProject(project *project.Project) *proto.Project {
	return &proto.Project{
		ProjectId:    project.ProjectId,
		ProjectCode:  project.ProjectCode,
		ProjectName:  project.ProjectName,
		ProjectDesc:  project.ProjectDesc,
		ProjectOwner: project.ProjectOwner,
		CreatedOn:    project.CreatedOn,
		CreatedBy:    project.CreatedBy,
		UpdatedOn:    project.UpdatedOn,
		UpdatedBy:    project.UpdatedBy,
		DeletedOn:    project.DeletedOn,
	}
}
