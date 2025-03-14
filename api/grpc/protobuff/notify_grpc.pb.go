// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: notify.proto

package protobuff

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NotificationService_HealthCheck_FullMethodName             = "/NotificationGRPC.NotificationService/HealthCheck"
	NotificationService_ListUsers_FullMethodName               = "/NotificationGRPC.NotificationService/ListUsers"
	NotificationService_RegisterSMS_FullMethodName             = "/NotificationGRPC.NotificationService/RegisterSMS"
	NotificationService_UnregisterSMS_FullMethodName           = "/NotificationGRPC.NotificationService/UnregisterSMS"
	NotificationService_UpdateSMSRegistration_FullMethodName   = "/NotificationGRPC.NotificationService/UpdateSMSRegistration"
	NotificationService_SMSNotify_FullMethodName               = "/NotificationGRPC.NotificationService/SMSNotify"
	NotificationService_SMSSendDirectMessage_FullMethodName    = "/NotificationGRPC.NotificationService/SMSSendDirectMessage"
	NotificationService_RegisterEmail_FullMethodName           = "/NotificationGRPC.NotificationService/RegisterEmail"
	NotificationService_UnregisterEmail_FullMethodName         = "/NotificationGRPC.NotificationService/UnregisterEmail"
	NotificationService_UpdateEmailRegistration_FullMethodName = "/NotificationGRPC.NotificationService/UpdateEmailRegistration"
	NotificationService_EmailNotify_FullMethodName             = "/NotificationGRPC.NotificationService/EmailNotify"
	NotificationService_EmailSendDirectMessage_FullMethodName  = "/NotificationGRPC.NotificationService/EmailSendDirectMessage"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	// Simple health check method
	HealthCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	// Service agnostic Methods
	ListUsers(ctx context.Context, in *Tags, opts ...grpc.CallOption) (*UserListResponse, error)
	// SMS methods
	// registration
	RegisterSMS(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	UnregisterSMS(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	UpdateSMSRegistration(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	// notification related
	SMSNotify(ctx context.Context, in *SMSNotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error)
	SMSSendDirectMessage(ctx context.Context, in *SMSSendDirectRequest, opts ...grpc.CallOption) (*ErrorArray, error)
	// Email Methods
	RegisterEmail(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	UnregisterEmail(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	UpdateEmailRegistration(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error)
	// notification related
	EmailNotify(ctx context.Context, in *EmailNotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error)
	EmailSendDirectMessage(ctx context.Context, in *EmailSendDirectRequest, opts ...grpc.CallOption) (*ErrorArray, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) HealthCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, NotificationService_HealthCheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) ListUsers(ctx context.Context, in *Tags, opts ...grpc.CallOption) (*UserListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserListResponse)
	err := c.cc.Invoke(ctx, NotificationService_ListUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) RegisterSMS(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_RegisterSMS_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UnregisterSMS(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_UnregisterSMS_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UpdateSMSRegistration(ctx context.Context, in *SMSRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_UpdateSMSRegistration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SMSNotify(ctx context.Context, in *SMSNotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NotifyResponse)
	err := c.cc.Invoke(ctx, NotificationService_SMSNotify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SMSSendDirectMessage(ctx context.Context, in *SMSSendDirectRequest, opts ...grpc.CallOption) (*ErrorArray, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ErrorArray)
	err := c.cc.Invoke(ctx, NotificationService_SMSSendDirectMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) RegisterEmail(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_RegisterEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UnregisterEmail(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_UnregisterEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) UpdateEmailRegistration(ctx context.Context, in *EmailRegisterInfo, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, NotificationService_UpdateEmailRegistration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) EmailNotify(ctx context.Context, in *EmailNotifyRequest, opts ...grpc.CallOption) (*NotifyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NotifyResponse)
	err := c.cc.Invoke(ctx, NotificationService_EmailNotify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) EmailSendDirectMessage(ctx context.Context, in *EmailSendDirectRequest, opts ...grpc.CallOption) (*ErrorArray, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ErrorArray)
	err := c.cc.Invoke(ctx, NotificationService_EmailSendDirectMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	// Simple health check method
	HealthCheck(context.Context, *empty.Empty) (*HealthCheckResponse, error)
	// Service agnostic Methods
	ListUsers(context.Context, *Tags) (*UserListResponse, error)
	// SMS methods
	// registration
	RegisterSMS(context.Context, *SMSRegisterInfo) (*BasicResponse, error)
	UnregisterSMS(context.Context, *SMSRegisterInfo) (*BasicResponse, error)
	UpdateSMSRegistration(context.Context, *SMSRegisterInfo) (*BasicResponse, error)
	// notification related
	SMSNotify(context.Context, *SMSNotifyRequest) (*NotifyResponse, error)
	SMSSendDirectMessage(context.Context, *SMSSendDirectRequest) (*ErrorArray, error)
	// Email Methods
	RegisterEmail(context.Context, *EmailRegisterInfo) (*BasicResponse, error)
	UnregisterEmail(context.Context, *EmailRegisterInfo) (*BasicResponse, error)
	UpdateEmailRegistration(context.Context, *EmailRegisterInfo) (*BasicResponse, error)
	// notification related
	EmailNotify(context.Context, *EmailNotifyRequest) (*NotifyResponse, error)
	EmailSendDirectMessage(context.Context, *EmailSendDirectRequest) (*ErrorArray, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) HealthCheck(context.Context, *empty.Empty) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedNotificationServiceServer) ListUsers(context.Context, *Tags) (*UserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedNotificationServiceServer) RegisterSMS(context.Context, *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterSMS not implemented")
}
func (UnimplementedNotificationServiceServer) UnregisterSMS(context.Context, *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterSMS not implemented")
}
func (UnimplementedNotificationServiceServer) UpdateSMSRegistration(context.Context, *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSMSRegistration not implemented")
}
func (UnimplementedNotificationServiceServer) SMSNotify(context.Context, *SMSNotifyRequest) (*NotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SMSNotify not implemented")
}
func (UnimplementedNotificationServiceServer) SMSSendDirectMessage(context.Context, *SMSSendDirectRequest) (*ErrorArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SMSSendDirectMessage not implemented")
}
func (UnimplementedNotificationServiceServer) RegisterEmail(context.Context, *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterEmail not implemented")
}
func (UnimplementedNotificationServiceServer) UnregisterEmail(context.Context, *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterEmail not implemented")
}
func (UnimplementedNotificationServiceServer) UpdateEmailRegistration(context.Context, *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEmailRegistration not implemented")
}
func (UnimplementedNotificationServiceServer) EmailNotify(context.Context, *EmailNotifyRequest) (*NotifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailNotify not implemented")
}
func (UnimplementedNotificationServiceServer) EmailSendDirectMessage(context.Context, *EmailSendDirectRequest) (*ErrorArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmailSendDirectMessage not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).HealthCheck(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tags)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_ListUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).ListUsers(ctx, req.(*Tags))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_RegisterSMS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).RegisterSMS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_RegisterSMS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).RegisterSMS(ctx, req.(*SMSRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UnregisterSMS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UnregisterSMS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UnregisterSMS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UnregisterSMS(ctx, req.(*SMSRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UpdateSMSRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UpdateSMSRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UpdateSMSRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UpdateSMSRegistration(ctx, req.(*SMSRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SMSNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSNotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SMSNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SMSNotify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SMSNotify(ctx, req.(*SMSNotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SMSSendDirectMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SMSSendDirectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SMSSendDirectMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SMSSendDirectMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SMSSendDirectMessage(ctx, req.(*SMSSendDirectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_RegisterEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).RegisterEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_RegisterEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).RegisterEmail(ctx, req.(*EmailRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UnregisterEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UnregisterEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UnregisterEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UnregisterEmail(ctx, req.(*EmailRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_UpdateEmailRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).UpdateEmailRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_UpdateEmailRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).UpdateEmailRegistration(ctx, req.(*EmailRegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_EmailNotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailNotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).EmailNotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_EmailNotify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).EmailNotify(ctx, req.(*EmailNotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_EmailSendDirectMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailSendDirectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).EmailSendDirectMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_EmailSendDirectMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).EmailSendDirectMessage(ctx, req.(*EmailSendDirectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NotificationGRPC.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _NotificationService_HealthCheck_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _NotificationService_ListUsers_Handler,
		},
		{
			MethodName: "RegisterSMS",
			Handler:    _NotificationService_RegisterSMS_Handler,
		},
		{
			MethodName: "UnregisterSMS",
			Handler:    _NotificationService_UnregisterSMS_Handler,
		},
		{
			MethodName: "UpdateSMSRegistration",
			Handler:    _NotificationService_UpdateSMSRegistration_Handler,
		},
		{
			MethodName: "SMSNotify",
			Handler:    _NotificationService_SMSNotify_Handler,
		},
		{
			MethodName: "SMSSendDirectMessage",
			Handler:    _NotificationService_SMSSendDirectMessage_Handler,
		},
		{
			MethodName: "RegisterEmail",
			Handler:    _NotificationService_RegisterEmail_Handler,
		},
		{
			MethodName: "UnregisterEmail",
			Handler:    _NotificationService_UnregisterEmail_Handler,
		},
		{
			MethodName: "UpdateEmailRegistration",
			Handler:    _NotificationService_UpdateEmailRegistration_Handler,
		},
		{
			MethodName: "EmailNotify",
			Handler:    _NotificationService_EmailNotify_Handler,
		},
		{
			MethodName: "EmailSendDirectMessage",
			Handler:    _NotificationService_EmailSendDirectMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notify.proto",
}
