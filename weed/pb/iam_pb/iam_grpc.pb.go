// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: iam.proto

package iam_pb

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const ()

// SeaweedIdentityAccessManagementClient is the client API for SeaweedIdentityAccessManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeaweedIdentityAccessManagementClient interface {
}

type seaweedIdentityAccessManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewSeaweedIdentityAccessManagementClient(cc grpc.ClientConnInterface) SeaweedIdentityAccessManagementClient {
	return &seaweedIdentityAccessManagementClient{cc}
}

// SeaweedIdentityAccessManagementServer is the server API for SeaweedIdentityAccessManagement service.
// All implementations must embed UnimplementedSeaweedIdentityAccessManagementServer
// for forward compatibility
type SeaweedIdentityAccessManagementServer interface {
	mustEmbedUnimplementedSeaweedIdentityAccessManagementServer()
}

// UnimplementedSeaweedIdentityAccessManagementServer must be embedded to have forward compatible implementations.
type UnimplementedSeaweedIdentityAccessManagementServer struct {
}

func (UnimplementedSeaweedIdentityAccessManagementServer) mustEmbedUnimplementedSeaweedIdentityAccessManagementServer() {
}

// UnsafeSeaweedIdentityAccessManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeaweedIdentityAccessManagementServer will
// result in compilation errors.
type UnsafeSeaweedIdentityAccessManagementServer interface {
	mustEmbedUnimplementedSeaweedIdentityAccessManagementServer()
}

func RegisterSeaweedIdentityAccessManagementServer(s grpc.ServiceRegistrar, srv SeaweedIdentityAccessManagementServer) {
	s.RegisterService(&SeaweedIdentityAccessManagement_ServiceDesc, srv)
}

// SeaweedIdentityAccessManagement_ServiceDesc is the grpc.ServiceDesc for SeaweedIdentityAccessManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SeaweedIdentityAccessManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iam_pb.SeaweedIdentityAccessManagement",
	HandlerType: (*SeaweedIdentityAccessManagementServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "iam.proto",
}
