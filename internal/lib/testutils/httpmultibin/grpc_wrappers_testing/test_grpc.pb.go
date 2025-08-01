// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: test.proto

package grpc_wrappers_testing

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Service_TestString_FullMethodName  = "/grpc.wrappers.testing.Service/TestString"
	Service_TestInteger_FullMethodName = "/grpc.wrappers.testing.Service/TestInteger"
	Service_TestBoolean_FullMethodName = "/grpc.wrappers.testing.Service/TestBoolean"
	Service_TestDouble_FullMethodName  = "/grpc.wrappers.testing.Service/TestDouble"
	Service_TestValue_FullMethodName   = "/grpc.wrappers.testing.Service/TestValue"
	Service_TestStream_FullMethodName  = "/grpc.wrappers.testing.Service/TestStream"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	TestString(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
	TestInteger(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*wrapperspb.Int64Value, error)
	TestBoolean(ctx context.Context, in *wrapperspb.BoolValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
	TestDouble(ctx context.Context, in *wrapperspb.DoubleValue, opts ...grpc.CallOption) (*wrapperspb.DoubleValue, error)
	TestValue(ctx context.Context, in *structpb.Value, opts ...grpc.CallOption) (*structpb.Value, error)
	TestStream(ctx context.Context, opts ...grpc.CallOption) (Service_TestStreamClient, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) TestString(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, Service_TestString_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) TestInteger(ctx context.Context, in *wrapperspb.Int64Value, opts ...grpc.CallOption) (*wrapperspb.Int64Value, error) {
	out := new(wrapperspb.Int64Value)
	err := c.cc.Invoke(ctx, Service_TestInteger_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) TestBoolean(ctx context.Context, in *wrapperspb.BoolValue, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, Service_TestBoolean_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) TestDouble(ctx context.Context, in *wrapperspb.DoubleValue, opts ...grpc.CallOption) (*wrapperspb.DoubleValue, error) {
	out := new(wrapperspb.DoubleValue)
	err := c.cc.Invoke(ctx, Service_TestDouble_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) TestValue(ctx context.Context, in *structpb.Value, opts ...grpc.CallOption) (*structpb.Value, error) {
	out := new(structpb.Value)
	err := c.cc.Invoke(ctx, Service_TestValue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) TestStream(ctx context.Context, opts ...grpc.CallOption) (Service_TestStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[0], Service_TestStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceTestStreamClient{stream}
	return x, nil
}

type Service_TestStreamClient interface {
	Send(*wrapperspb.StringValue) error
	CloseAndRecv() (*wrapperspb.StringValue, error)
	grpc.ClientStream
}

type serviceTestStreamClient struct {
	grpc.ClientStream
}

func (x *serviceTestStreamClient) Send(m *wrapperspb.StringValue) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceTestStreamClient) CloseAndRecv() (*wrapperspb.StringValue, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(wrapperspb.StringValue)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	TestString(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error)
	TestInteger(context.Context, *wrapperspb.Int64Value) (*wrapperspb.Int64Value, error)
	TestBoolean(context.Context, *wrapperspb.BoolValue) (*wrapperspb.BoolValue, error)
	TestDouble(context.Context, *wrapperspb.DoubleValue) (*wrapperspb.DoubleValue, error)
	TestValue(context.Context, *structpb.Value) (*structpb.Value, error)
	TestStream(Service_TestStreamServer) error
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) TestString(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestString not implemented")
}
func (UnimplementedServiceServer) TestInteger(context.Context, *wrapperspb.Int64Value) (*wrapperspb.Int64Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestInteger not implemented")
}
func (UnimplementedServiceServer) TestBoolean(context.Context, *wrapperspb.BoolValue) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestBoolean not implemented")
}
func (UnimplementedServiceServer) TestDouble(context.Context, *wrapperspb.DoubleValue) (*wrapperspb.DoubleValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestDouble not implemented")
}
func (UnimplementedServiceServer) TestValue(context.Context, *structpb.Value) (*structpb.Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestValue not implemented")
}
func (UnimplementedServiceServer) TestStream(Service_TestStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method TestStream not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_TestString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).TestString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_TestString_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).TestString(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_TestInteger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.Int64Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).TestInteger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_TestInteger_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).TestInteger(ctx, req.(*wrapperspb.Int64Value))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_TestBoolean_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.BoolValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).TestBoolean(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_TestBoolean_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).TestBoolean(ctx, req.(*wrapperspb.BoolValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_TestDouble_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.DoubleValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).TestDouble(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_TestDouble_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).TestDouble(ctx, req.(*wrapperspb.DoubleValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_TestValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(structpb.Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).TestValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_TestValue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).TestValue(ctx, req.(*structpb.Value))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_TestStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).TestStream(&serviceTestStreamServer{stream})
}

type Service_TestStreamServer interface {
	SendAndClose(*wrapperspb.StringValue) error
	Recv() (*wrapperspb.StringValue, error)
	grpc.ServerStream
}

type serviceTestStreamServer struct {
	grpc.ServerStream
}

func (x *serviceTestStreamServer) SendAndClose(m *wrapperspb.StringValue) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceTestStreamServer) Recv() (*wrapperspb.StringValue, error) {
	m := new(wrapperspb.StringValue)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.wrappers.testing.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TestString",
			Handler:    _Service_TestString_Handler,
		},
		{
			MethodName: "TestInteger",
			Handler:    _Service_TestInteger_Handler,
		},
		{
			MethodName: "TestBoolean",
			Handler:    _Service_TestBoolean_Handler,
		},
		{
			MethodName: "TestDouble",
			Handler:    _Service_TestDouble_Handler,
		},
		{
			MethodName: "TestValue",
			Handler:    _Service_TestValue_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TestStream",
			Handler:       _Service_TestStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "test.proto",
}
