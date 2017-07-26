// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	GetProjectsResponse
	GetBuildsRequest
	GetBuildsResponse
	NewBuildRequest
	NextJobRequest
	NextJobResponse
	JobDoneRequest
	Build
	Project
	StandardRequest
	StandardResponse
	RequestHeader
	ResponseHeader
	Error
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetProjectsResponse struct {
	Header   *ResponseHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Projects []string        `protobuf:"bytes,2,rep,name=projects" json:"projects,omitempty"`
}

func (m *GetProjectsResponse) Reset()                    { *m = GetProjectsResponse{} }
func (m *GetProjectsResponse) String() string            { return proto.CompactTextString(m) }
func (*GetProjectsResponse) ProtoMessage()               {}
func (*GetProjectsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetProjectsResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetProjectsResponse) GetProjects() []string {
	if m != nil {
		return m.Projects
	}
	return nil
}

type GetBuildsRequest struct {
	Header      *RequestHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	ProjectName string         `protobuf:"bytes,2,opt,name=projectName" json:"projectName,omitempty"`
}

func (m *GetBuildsRequest) Reset()                    { *m = GetBuildsRequest{} }
func (m *GetBuildsRequest) String() string            { return proto.CompactTextString(m) }
func (*GetBuildsRequest) ProtoMessage()               {}
func (*GetBuildsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetBuildsRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetBuildsRequest) GetProjectName() string {
	if m != nil {
		return m.ProjectName
	}
	return ""
}

type GetBuildsResponse struct {
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Builds []*Build        `protobuf:"bytes,2,rep,name=builds" json:"builds,omitempty"`
}

func (m *GetBuildsResponse) Reset()                    { *m = GetBuildsResponse{} }
func (m *GetBuildsResponse) String() string            { return proto.CompactTextString(m) }
func (*GetBuildsResponse) ProtoMessage()               {}
func (*GetBuildsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetBuildsResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetBuildsResponse) GetBuilds() []*Build {
	if m != nil {
		return m.Builds
	}
	return nil
}

type NewBuildRequest struct {
	Header *RequestHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Build  *Build         `protobuf:"bytes,2,opt,name=build" json:"build,omitempty"`
}

func (m *NewBuildRequest) Reset()                    { *m = NewBuildRequest{} }
func (m *NewBuildRequest) String() string            { return proto.CompactTextString(m) }
func (*NewBuildRequest) ProtoMessage()               {}
func (*NewBuildRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NewBuildRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *NewBuildRequest) GetBuild() *Build {
	if m != nil {
		return m.Build
	}
	return nil
}

type NextJobRequest struct {
	Header *RequestHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *NextJobRequest) Reset()                    { *m = NextJobRequest{} }
func (m *NextJobRequest) String() string            { return proto.CompactTextString(m) }
func (*NextJobRequest) ProtoMessage()               {}
func (*NextJobRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *NextJobRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type NextJobResponse struct {
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *NextJobResponse) Reset()                    { *m = NextJobResponse{} }
func (m *NextJobResponse) String() string            { return proto.CompactTextString(m) }
func (*NextJobResponse) ProtoMessage()               {}
func (*NextJobResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *NextJobResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type JobDoneRequest struct {
	Header *RequestHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *JobDoneRequest) Reset()                    { *m = JobDoneRequest{} }
func (m *JobDoneRequest) String() string            { return proto.CompactTextString(m) }
func (*JobDoneRequest) ProtoMessage()               {}
func (*JobDoneRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *JobDoneRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type Build struct {
	ProjectName    string `protobuf:"bytes,1,opt,name=projectName" json:"projectName,omitempty"`
	BuildId        string `protobuf:"bytes,2,opt,name=buildId" json:"buildId,omitempty"`
	BuildURL       string `protobuf:"bytes,3,opt,name=buildURL" json:"buildURL,omitempty"`
	BuildSystemURL string `protobuf:"bytes,4,opt,name=buildSystemURL" json:"buildSystemURL,omitempty"`
	BranchName     string `protobuf:"bytes,5,opt,name=branchName" json:"branchName,omitempty"`
	CreatedAt      int64  `protobuf:"varint,6,opt,name=createdAt" json:"createdAt,omitempty"`
}

func (m *Build) Reset()                    { *m = Build{} }
func (m *Build) String() string            { return proto.CompactTextString(m) }
func (*Build) ProtoMessage()               {}
func (*Build) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Build) GetProjectName() string {
	if m != nil {
		return m.ProjectName
	}
	return ""
}

func (m *Build) GetBuildId() string {
	if m != nil {
		return m.BuildId
	}
	return ""
}

func (m *Build) GetBuildURL() string {
	if m != nil {
		return m.BuildURL
	}
	return ""
}

func (m *Build) GetBuildSystemURL() string {
	if m != nil {
		return m.BuildSystemURL
	}
	return ""
}

func (m *Build) GetBranchName() string {
	if m != nil {
		return m.BranchName
	}
	return ""
}

func (m *Build) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

type Project struct {
	ProjectName string   `protobuf:"bytes,1,opt,name=projectName" json:"projectName,omitempty"`
	Builds      []*Build `protobuf:"bytes,2,rep,name=builds" json:"builds,omitempty"`
}

func (m *Project) Reset()                    { *m = Project{} }
func (m *Project) String() string            { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()               {}
func (*Project) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Project) GetProjectName() string {
	if m != nil {
		return m.ProjectName
	}
	return ""
}

func (m *Project) GetBuilds() []*Build {
	if m != nil {
		return m.Builds
	}
	return nil
}

type StandardRequest struct {
	Header *RequestHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *StandardRequest) Reset()                    { *m = StandardRequest{} }
func (m *StandardRequest) String() string            { return proto.CompactTextString(m) }
func (*StandardRequest) ProtoMessage()               {}
func (*StandardRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *StandardRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type StandardResponse struct {
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *StandardResponse) Reset()                    { *m = StandardResponse{} }
func (m *StandardResponse) String() string            { return proto.CompactTextString(m) }
func (*StandardResponse) ProtoMessage()               {}
func (*StandardResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *StandardResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type RequestHeader struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *RequestHeader) Reset()                    { *m = RequestHeader{} }
func (m *RequestHeader) String() string            { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()               {}
func (*RequestHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *RequestHeader) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type ResponseHeader struct {
	Success bool     `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Errors  []*Error `protobuf:"bytes,2,rep,name=errors" json:"errors,omitempty"`
}

func (m *ResponseHeader) Reset()                    { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string            { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()               {}
func (*ResponseHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *ResponseHeader) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ResponseHeader) GetErrors() []*Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

type Error struct {
	Code    string `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Error) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*GetProjectsResponse)(nil), "protobuf.GetProjectsResponse")
	proto.RegisterType((*GetBuildsRequest)(nil), "protobuf.GetBuildsRequest")
	proto.RegisterType((*GetBuildsResponse)(nil), "protobuf.GetBuildsResponse")
	proto.RegisterType((*NewBuildRequest)(nil), "protobuf.NewBuildRequest")
	proto.RegisterType((*NextJobRequest)(nil), "protobuf.NextJobRequest")
	proto.RegisterType((*NextJobResponse)(nil), "protobuf.NextJobResponse")
	proto.RegisterType((*JobDoneRequest)(nil), "protobuf.JobDoneRequest")
	proto.RegisterType((*Build)(nil), "protobuf.Build")
	proto.RegisterType((*Project)(nil), "protobuf.Project")
	proto.RegisterType((*StandardRequest)(nil), "protobuf.StandardRequest")
	proto.RegisterType((*StandardResponse)(nil), "protobuf.StandardResponse")
	proto.RegisterType((*RequestHeader)(nil), "protobuf.RequestHeader")
	proto.RegisterType((*ResponseHeader)(nil), "protobuf.ResponseHeader")
	proto.RegisterType((*Error)(nil), "protobuf.Error")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DeploiServer service

type DeploiServerClient interface {
	RegisterNewBuild(ctx context.Context, in *NewBuildRequest, opts ...grpc.CallOption) (*StandardResponse, error)
	GetNextJob(ctx context.Context, in *NextJobRequest, opts ...grpc.CallOption) (*NextJobResponse, error)
	MarkJobDone(ctx context.Context, in *JobDoneRequest, opts ...grpc.CallOption) (*StandardResponse, error)
	GetProjects(ctx context.Context, in *StandardRequest, opts ...grpc.CallOption) (*GetProjectsResponse, error)
	GetBuilds(ctx context.Context, in *GetBuildsRequest, opts ...grpc.CallOption) (*GetBuildsResponse, error)
}

type deploiServerClient struct {
	cc *grpc.ClientConn
}

func NewDeploiServerClient(cc *grpc.ClientConn) DeploiServerClient {
	return &deploiServerClient{cc}
}

func (c *deploiServerClient) RegisterNewBuild(ctx context.Context, in *NewBuildRequest, opts ...grpc.CallOption) (*StandardResponse, error) {
	out := new(StandardResponse)
	err := grpc.Invoke(ctx, "/protobuf.DeploiServer/RegisterNewBuild", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deploiServerClient) GetNextJob(ctx context.Context, in *NextJobRequest, opts ...grpc.CallOption) (*NextJobResponse, error) {
	out := new(NextJobResponse)
	err := grpc.Invoke(ctx, "/protobuf.DeploiServer/GetNextJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deploiServerClient) MarkJobDone(ctx context.Context, in *JobDoneRequest, opts ...grpc.CallOption) (*StandardResponse, error) {
	out := new(StandardResponse)
	err := grpc.Invoke(ctx, "/protobuf.DeploiServer/MarkJobDone", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deploiServerClient) GetProjects(ctx context.Context, in *StandardRequest, opts ...grpc.CallOption) (*GetProjectsResponse, error) {
	out := new(GetProjectsResponse)
	err := grpc.Invoke(ctx, "/protobuf.DeploiServer/GetProjects", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deploiServerClient) GetBuilds(ctx context.Context, in *GetBuildsRequest, opts ...grpc.CallOption) (*GetBuildsResponse, error) {
	out := new(GetBuildsResponse)
	err := grpc.Invoke(ctx, "/protobuf.DeploiServer/GetBuilds", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DeploiServer service

type DeploiServerServer interface {
	RegisterNewBuild(context.Context, *NewBuildRequest) (*StandardResponse, error)
	GetNextJob(context.Context, *NextJobRequest) (*NextJobResponse, error)
	MarkJobDone(context.Context, *JobDoneRequest) (*StandardResponse, error)
	GetProjects(context.Context, *StandardRequest) (*GetProjectsResponse, error)
	GetBuilds(context.Context, *GetBuildsRequest) (*GetBuildsResponse, error)
}

func RegisterDeploiServerServer(s *grpc.Server, srv DeploiServerServer) {
	s.RegisterService(&_DeploiServer_serviceDesc, srv)
}

func _DeploiServer_RegisterNewBuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewBuildRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploiServerServer).RegisterNewBuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DeploiServer/RegisterNewBuild",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploiServerServer).RegisterNewBuild(ctx, req.(*NewBuildRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploiServer_GetNextJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NextJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploiServerServer).GetNextJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DeploiServer/GetNextJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploiServerServer).GetNextJob(ctx, req.(*NextJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploiServer_MarkJobDone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobDoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploiServerServer).MarkJobDone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DeploiServer/MarkJobDone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploiServerServer).MarkJobDone(ctx, req.(*JobDoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploiServer_GetProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StandardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploiServerServer).GetProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DeploiServer/GetProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploiServerServer).GetProjects(ctx, req.(*StandardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploiServer_GetBuilds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBuildsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploiServerServer).GetBuilds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.DeploiServer/GetBuilds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploiServerServer).GetBuilds(ctx, req.(*GetBuildsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DeploiServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.DeploiServer",
	HandlerType: (*DeploiServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNewBuild",
			Handler:    _DeploiServer_RegisterNewBuild_Handler,
		},
		{
			MethodName: "GetNextJob",
			Handler:    _DeploiServer_GetNextJob_Handler,
		},
		{
			MethodName: "MarkJobDone",
			Handler:    _DeploiServer_MarkJobDone_Handler,
		},
		{
			MethodName: "GetProjects",
			Handler:    _DeploiServer_GetProjects_Handler,
		},
		{
			MethodName: "GetBuilds",
			Handler:    _DeploiServer_GetBuilds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0x93, 0x26, 0x6d, 0x26, 0x25, 0x09, 0x0b, 0x12, 0xae, 0xb9, 0x28, 0x5a, 0xa9, 0x90,
	0xa7, 0x80, 0x82, 0xf8, 0x80, 0xb6, 0x29, 0xa1, 0x05, 0x22, 0xb4, 0x81, 0x0f, 0xf0, 0x65, 0x68,
	0x43, 0x1b, 0x6f, 0xd8, 0xdd, 0x70, 0xf9, 0x40, 0x7e, 0x8a, 0x27, 0xe4, 0xc9, 0x3a, 0x76, 0x8c,
	0xa5, 0x56, 0x7e, 0xf2, 0xce, 0x9c, 0xf1, 0x99, 0x3d, 0x33, 0x3b, 0x03, 0x07, 0x1a, 0xd5, 0x0f,
	0x54, 0xc3, 0xa5, 0x92, 0x46, 0xb2, 0x7d, 0xfa, 0x04, 0xab, 0xaf, 0x3c, 0x84, 0x07, 0x13, 0x34,
	0x9f, 0x94, 0xfc, 0x86, 0xa1, 0xd1, 0x02, 0xf5, 0x52, 0xc6, 0x1a, 0xd9, 0x2b, 0x68, 0x5e, 0xa1,
	0x1f, 0xa1, 0x72, 0x9d, 0xbe, 0x33, 0x68, 0x8f, 0xdc, 0x61, 0xfa, 0xc7, 0x30, 0x8d, 0x79, 0x47,
	0xb8, 0xb0, 0x71, 0xcc, 0x83, 0x84, 0x94, 0x58, 0xdc, 0x5a, 0xbf, 0x3e, 0x68, 0x89, 0x8d, 0xcd,
	0x11, 0x7a, 0x13, 0x34, 0x27, 0xab, 0xf9, 0x4d, 0xa4, 0x05, 0x7e, 0x5f, 0xa1, 0x36, 0xec, 0x65,
	0x21, 0xc3, 0xa3, 0x7c, 0x06, 0x0a, 0x29, 0x24, 0xe8, 0x43, 0xdb, 0x12, 0x4e, 0xfd, 0x05, 0xba,
	0xb5, 0xbe, 0x33, 0x68, 0x89, 0xbc, 0x8b, 0xc7, 0x70, 0x3f, 0x97, 0xa6, 0xb2, 0x92, 0x17, 0xd0,
	0x0c, 0x88, 0x83, 0x74, 0xb4, 0x47, 0xdd, 0xec, 0x0f, 0xe2, 0x16, 0x16, 0xe6, 0x73, 0xe8, 0x4e,
	0xf1, 0xe7, 0xda, 0x57, 0x55, 0xd5, 0x11, 0x34, 0x88, 0x8d, 0xf4, 0x94, 0xe4, 0x5a, 0xa3, 0xfc,
	0x18, 0x3a, 0x53, 0xfc, 0x65, 0x2e, 0x64, 0x50, 0x35, 0x13, 0x3f, 0x4d, 0x6e, 0x6b, 0x29, 0xaa,
	0xd6, 0x26, 0xb9, 0xc7, 0x85, 0x0c, 0xc6, 0x32, 0xc6, 0xca, 0xf7, 0xf8, 0xe3, 0x40, 0x83, 0xb4,
	0x15, 0x3b, 0xea, 0xfc, 0xd7, 0x51, 0xe6, 0xc2, 0x1e, 0xe9, 0x3f, 0x8f, 0x6c, 0xbf, 0x53, 0x33,
	0x79, 0x6e, 0x74, 0xfc, 0x22, 0x3e, 0xb8, 0x75, 0x82, 0x36, 0x36, 0x7b, 0x0e, 0x1d, 0x3a, 0xcf,
	0x7e, 0x6b, 0x83, 0x8b, 0x24, 0x62, 0x97, 0x22, 0x0a, 0x5e, 0xf6, 0x0c, 0x20, 0x50, 0x7e, 0x1c,
	0x5e, 0x51, 0xfa, 0x06, 0xc5, 0xe4, 0x3c, 0xec, 0x09, 0xb4, 0x42, 0x85, 0xbe, 0xc1, 0xe8, 0xd8,
	0xb8, 0xcd, 0xbe, 0x33, 0xa8, 0x8b, 0xcc, 0xc1, 0x3f, 0xc3, 0x9e, 0x1d, 0x9b, 0x3b, 0x08, 0xb9,
	0xf3, 0x9b, 0x3a, 0x81, 0xee, 0xcc, 0xf8, 0x71, 0xe4, 0xab, 0xca, 0x6f, 0x8a, 0x8f, 0xa1, 0x97,
	0x71, 0x54, 0x6e, 0xf5, 0x11, 0xdc, 0xdb, 0xa2, 0x67, 0x0f, 0xa1, 0x61, 0xe4, 0x35, 0xc6, 0x56,
	0xdf, 0xda, 0xe0, 0x33, 0xe8, 0x6c, 0x13, 0x24, 0x4d, 0xd3, 0xab, 0x30, 0x44, 0xad, 0x29, 0x72,
	0x5f, 0xa4, 0x66, 0x52, 0x05, 0x54, 0x4a, 0xaa, 0x92, 0x2a, 0x9c, 0x25, 0x7e, 0x61, 0x61, 0xfe,
	0x06, 0x1a, 0xe4, 0x60, 0x0c, 0x76, 0x43, 0x19, 0xa5, 0x25, 0xa5, 0x73, 0xc2, 0xbf, 0x40, 0xad,
	0xfd, 0xcb, 0x74, 0x09, 0xa4, 0xe6, 0xe8, 0x6f, 0x0d, 0x0e, 0xc6, 0xb8, 0xbc, 0x91, 0xf3, 0x19,
	0x6d, 0x3b, 0xf6, 0x1e, 0x7a, 0x02, 0x2f, 0xe7, 0xda, 0xa0, 0x4a, 0x27, 0x95, 0x1d, 0x66, 0x49,
	0x0b, 0xd3, 0xeb, 0x79, 0x19, 0x54, 0x2c, 0x20, 0xdf, 0x61, 0xa7, 0x00, 0x13, 0x34, 0x76, 0x86,
	0x98, 0x9b, 0xa7, 0xc9, 0x4f, 0xa6, 0x77, 0x58, 0x82, 0x6c, 0x48, 0xce, 0xa0, 0xfd, 0xd1, 0x57,
	0xd7, 0x76, 0x88, 0xf2, 0x2c, 0xdb, 0x73, 0x75, 0xcb, 0x5d, 0xce, 0xa1, 0x9d, 0x5b, 0xdb, 0x79,
	0x4d, 0x85, 0xd7, 0xe3, 0x3d, 0xcd, 0xa0, 0x92, 0x45, 0xcf, 0x77, 0xd8, 0x5b, 0x68, 0x6d, 0xb6,
	0x26, 0xf3, 0xb6, 0xa2, 0xb7, 0x36, 0xb6, 0xf7, 0xb8, 0x14, 0x4b, 0x79, 0x82, 0x26, 0xa1, 0xaf,
	0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x41, 0x37, 0x23, 0x39, 0x6a, 0x06, 0x00, 0x00,
}
