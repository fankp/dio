// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: project.proto

package devops_integral_upm_srv_service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ProjectService service

type ProjectService interface {
	// 创建项目
	CreateProject(ctx context.Context, in *CreateProjectReq, opts ...client.CallOption) (*ProjectResp, error)
	// 更新项目
	UpdateProject(ctx context.Context, in *UpdateProjectReq, opts ...client.CallOption) (*ProjectResp, error)
	// 查询用户关联的项目
	SelectUserProjects(ctx context.Context, in *SelectUserProjectsReq, opts ...client.CallOption) (*UserProjectsResp, error)
}

type projectService struct {
	c    client.Client
	name string
}

func NewProjectService(name string, c client.Client) ProjectService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "devops.integral.upm.srv.service"
	}
	return &projectService{
		c:    c,
		name: name,
	}
}

func (c *projectService) CreateProject(ctx context.Context, in *CreateProjectReq, opts ...client.CallOption) (*ProjectResp, error) {
	req := c.c.NewRequest(c.name, "ProjectService.CreateProject", in)
	out := new(ProjectResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectService) UpdateProject(ctx context.Context, in *UpdateProjectReq, opts ...client.CallOption) (*ProjectResp, error) {
	req := c.c.NewRequest(c.name, "ProjectService.UpdateProject", in)
	out := new(ProjectResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectService) SelectUserProjects(ctx context.Context, in *SelectUserProjectsReq, opts ...client.CallOption) (*UserProjectsResp, error) {
	req := c.c.NewRequest(c.name, "ProjectService.SelectUserProjects", in)
	out := new(UserProjectsResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProjectService service

type ProjectServiceHandler interface {
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq, *ProjectResp) error
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq, *ProjectResp) error
	// 查询用户关联的项目
	SelectUserProjects(context.Context, *SelectUserProjectsReq, *UserProjectsResp) error
}

func RegisterProjectServiceHandler(s server.Server, hdlr ProjectServiceHandler, opts ...server.HandlerOption) error {
	type projectService interface {
		CreateProject(ctx context.Context, in *CreateProjectReq, out *ProjectResp) error
		UpdateProject(ctx context.Context, in *UpdateProjectReq, out *ProjectResp) error
		SelectUserProjects(ctx context.Context, in *SelectUserProjectsReq, out *UserProjectsResp) error
	}
	type ProjectService struct {
		projectService
	}
	h := &projectServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ProjectService{h}, opts...))
}

type projectServiceHandler struct {
	ProjectServiceHandler
}

func (h *projectServiceHandler) CreateProject(ctx context.Context, in *CreateProjectReq, out *ProjectResp) error {
	return h.ProjectServiceHandler.CreateProject(ctx, in, out)
}

func (h *projectServiceHandler) UpdateProject(ctx context.Context, in *UpdateProjectReq, out *ProjectResp) error {
	return h.ProjectServiceHandler.UpdateProject(ctx, in, out)
}

func (h *projectServiceHandler) SelectUserProjects(ctx context.Context, in *SelectUserProjectsReq, out *UserProjectsResp) error {
	return h.ProjectServiceHandler.SelectUserProjects(ctx, in, out)
}