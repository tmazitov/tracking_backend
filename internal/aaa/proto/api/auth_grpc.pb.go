// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: files/auth.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AAAClient is the client API for AAA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AAAClient interface {
	CheckAuth(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
}

type aAAClient struct {
	cc grpc.ClientConnInterface
}

func NewAAAClient(cc grpc.ClientConnInterface) AAAClient {
	return &aAAClient{cc}
}

func (c *aAAClient) CheckAuth(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/auth.AAA/CheckAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AAAServer is the server API for AAA service.
// All implementations should embed UnimplementedAAAServer
// for forward compatibility
type AAAServer interface {
	CheckAuth(context.Context, *CheckRequest) (*CheckResponse, error)
}

// UnimplementedAAAServer should be embedded to have forward compatible implementations.
type UnimplementedAAAServer struct {
}

func (UnimplementedAAAServer) CheckAuth(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuth not implemented")
}

// UnsafeAAAServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AAAServer will
// result in compilation errors.
type UnsafeAAAServer interface {
	mustEmbedUnimplementedAAAServer()
}

func RegisterAAAServer(s grpc.ServiceRegistrar, srv AAAServer) {
	s.RegisterService(&AAA_ServiceDesc, srv)
}

func _AAA_CheckAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AAAServer).CheckAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AAA/CheckAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AAAServer).CheckAuth(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AAA_ServiceDesc is the grpc.ServiceDesc for AAA service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AAA_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AAA",
	HandlerType: (*AAAServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckAuth",
			Handler:    _AAA_CheckAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "files/auth.proto",
}
