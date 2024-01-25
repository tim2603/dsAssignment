// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: Communication.proto

package ds

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

// CommunicationWithMasterServiceClient is the client API for CommunicationWithMasterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationWithMasterServiceClient interface {
	SendAddressToMaster(ctx context.Context, in *Address, opts ...grpc.CallOption) (*EmptyMessage, error)
	SendHeartBeatToMaster(ctx context.Context, in *Heartbeat, opts ...grpc.CallOption) (*EmptyMessage, error)
	NotifyAboutFinishedMapTask(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*EmptyMessage, error)
	NotifyAboutFinishedReduceTask(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*EmptyMessage, error)
	SendCurrentTournamentTreeValue(ctx context.Context, in *TournamentTreeValue, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type communicationWithMasterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationWithMasterServiceClient(cc grpc.ClientConnInterface) CommunicationWithMasterServiceClient {
	return &communicationWithMasterServiceClient{cc}
}

func (c *communicationWithMasterServiceClient) SendAddressToMaster(ctx context.Context, in *Address, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithMasterService/SendAddressToMaster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithMasterServiceClient) SendHeartBeatToMaster(ctx context.Context, in *Heartbeat, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithMasterService/SendHeartBeatToMaster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithMasterServiceClient) NotifyAboutFinishedMapTask(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithMasterService/NotifyAboutFinishedMapTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithMasterServiceClient) NotifyAboutFinishedReduceTask(ctx context.Context, in *WorkerID, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithMasterService/NotifyAboutFinishedReduceTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithMasterServiceClient) SendCurrentTournamentTreeValue(ctx context.Context, in *TournamentTreeValue, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithMasterService/SendCurrentTournamentTreeValue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationWithMasterServiceServer is the server API for CommunicationWithMasterService service.
// All implementations must embed UnimplementedCommunicationWithMasterServiceServer
// for forward compatibility
type CommunicationWithMasterServiceServer interface {
	SendAddressToMaster(context.Context, *Address) (*EmptyMessage, error)
	SendHeartBeatToMaster(context.Context, *Heartbeat) (*EmptyMessage, error)
	NotifyAboutFinishedMapTask(context.Context, *WorkerID) (*EmptyMessage, error)
	NotifyAboutFinishedReduceTask(context.Context, *WorkerID) (*EmptyMessage, error)
	SendCurrentTournamentTreeValue(context.Context, *TournamentTreeValue) (*EmptyMessage, error)
	mustEmbedUnimplementedCommunicationWithMasterServiceServer()
}

// UnimplementedCommunicationWithMasterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicationWithMasterServiceServer struct {
}

func (UnimplementedCommunicationWithMasterServiceServer) SendAddressToMaster(context.Context, *Address) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendAddressToMaster not implemented")
}
func (UnimplementedCommunicationWithMasterServiceServer) SendHeartBeatToMaster(context.Context, *Heartbeat) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartBeatToMaster not implemented")
}
func (UnimplementedCommunicationWithMasterServiceServer) NotifyAboutFinishedMapTask(context.Context, *WorkerID) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyAboutFinishedMapTask not implemented")
}
func (UnimplementedCommunicationWithMasterServiceServer) NotifyAboutFinishedReduceTask(context.Context, *WorkerID) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyAboutFinishedReduceTask not implemented")
}
func (UnimplementedCommunicationWithMasterServiceServer) SendCurrentTournamentTreeValue(context.Context, *TournamentTreeValue) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCurrentTournamentTreeValue not implemented")
}
func (UnimplementedCommunicationWithMasterServiceServer) mustEmbedUnimplementedCommunicationWithMasterServiceServer() {
}

// UnsafeCommunicationWithMasterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationWithMasterServiceServer will
// result in compilation errors.
type UnsafeCommunicationWithMasterServiceServer interface {
	mustEmbedUnimplementedCommunicationWithMasterServiceServer()
}

func RegisterCommunicationWithMasterServiceServer(s grpc.ServiceRegistrar, srv CommunicationWithMasterServiceServer) {
	s.RegisterService(&CommunicationWithMasterService_ServiceDesc, srv)
}

func _CommunicationWithMasterService_SendAddressToMaster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithMasterServiceServer).SendAddressToMaster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithMasterService/SendAddressToMaster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithMasterServiceServer).SendAddressToMaster(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithMasterService_SendHeartBeatToMaster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Heartbeat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithMasterServiceServer).SendHeartBeatToMaster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithMasterService/SendHeartBeatToMaster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithMasterServiceServer).SendHeartBeatToMaster(ctx, req.(*Heartbeat))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithMasterService_NotifyAboutFinishedMapTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkerID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithMasterServiceServer).NotifyAboutFinishedMapTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithMasterService/NotifyAboutFinishedMapTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithMasterServiceServer).NotifyAboutFinishedMapTask(ctx, req.(*WorkerID))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithMasterService_NotifyAboutFinishedReduceTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkerID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithMasterServiceServer).NotifyAboutFinishedReduceTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithMasterService/NotifyAboutFinishedReduceTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithMasterServiceServer).NotifyAboutFinishedReduceTask(ctx, req.(*WorkerID))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithMasterService_SendCurrentTournamentTreeValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TournamentTreeValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithMasterServiceServer).SendCurrentTournamentTreeValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithMasterService/SendCurrentTournamentTreeValue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithMasterServiceServer).SendCurrentTournamentTreeValue(ctx, req.(*TournamentTreeValue))
	}
	return interceptor(ctx, in, info, handler)
}

// CommunicationWithMasterService_ServiceDesc is the grpc.ServiceDesc for CommunicationWithMasterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommunicationWithMasterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ds.CommunicationWithMasterService",
	HandlerType: (*CommunicationWithMasterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendAddressToMaster",
			Handler:    _CommunicationWithMasterService_SendAddressToMaster_Handler,
		},
		{
			MethodName: "SendHeartBeatToMaster",
			Handler:    _CommunicationWithMasterService_SendHeartBeatToMaster_Handler,
		},
		{
			MethodName: "NotifyAboutFinishedMapTask",
			Handler:    _CommunicationWithMasterService_NotifyAboutFinishedMapTask_Handler,
		},
		{
			MethodName: "NotifyAboutFinishedReduceTask",
			Handler:    _CommunicationWithMasterService_NotifyAboutFinishedReduceTask_Handler,
		},
		{
			MethodName: "SendCurrentTournamentTreeValue",
			Handler:    _CommunicationWithMasterService_SendCurrentTournamentTreeValue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Communication.proto",
}

// CommunicationWithWorkerServiceClient is the client API for CommunicationWithWorkerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationWithWorkerServiceClient interface {
	AssignMapTask(ctx context.Context, in *MapTask, opts ...grpc.CallOption) (*EmptyMessage, error)
	AssignReduceTask(ctx context.Context, in *ReduceTask, opts ...grpc.CallOption) (*EmptyMessage, error)
	SendSelectedTournamentTreeSelection(ctx context.Context, in *TournamentTreeSelection, opts ...grpc.CallOption) (*EmptyMessage, error)
	RestartTournamentTree(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type communicationWithWorkerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationWithWorkerServiceClient(cc grpc.ClientConnInterface) CommunicationWithWorkerServiceClient {
	return &communicationWithWorkerServiceClient{cc}
}

func (c *communicationWithWorkerServiceClient) AssignMapTask(ctx context.Context, in *MapTask, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithWorkerService/assignMapTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithWorkerServiceClient) AssignReduceTask(ctx context.Context, in *ReduceTask, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithWorkerService/assignReduceTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithWorkerServiceClient) SendSelectedTournamentTreeSelection(ctx context.Context, in *TournamentTreeSelection, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithWorkerService/sendSelectedTournamentTreeSelection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationWithWorkerServiceClient) RestartTournamentTree(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ds.CommunicationWithWorkerService/restartTournamentTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationWithWorkerServiceServer is the server API for CommunicationWithWorkerService service.
// All implementations must embed UnimplementedCommunicationWithWorkerServiceServer
// for forward compatibility
type CommunicationWithWorkerServiceServer interface {
	AssignMapTask(context.Context, *MapTask) (*EmptyMessage, error)
	AssignReduceTask(context.Context, *ReduceTask) (*EmptyMessage, error)
	SendSelectedTournamentTreeSelection(context.Context, *TournamentTreeSelection) (*EmptyMessage, error)
	RestartTournamentTree(context.Context, *EmptyMessage) (*EmptyMessage, error)
	mustEmbedUnimplementedCommunicationWithWorkerServiceServer()
}

// UnimplementedCommunicationWithWorkerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicationWithWorkerServiceServer struct {
}

func (UnimplementedCommunicationWithWorkerServiceServer) AssignMapTask(context.Context, *MapTask) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignMapTask not implemented")
}
func (UnimplementedCommunicationWithWorkerServiceServer) AssignReduceTask(context.Context, *ReduceTask) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignReduceTask not implemented")
}
func (UnimplementedCommunicationWithWorkerServiceServer) SendSelectedTournamentTreeSelection(context.Context, *TournamentTreeSelection) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSelectedTournamentTreeSelection not implemented")
}
func (UnimplementedCommunicationWithWorkerServiceServer) RestartTournamentTree(context.Context, *EmptyMessage) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartTournamentTree not implemented")
}
func (UnimplementedCommunicationWithWorkerServiceServer) mustEmbedUnimplementedCommunicationWithWorkerServiceServer() {
}

// UnsafeCommunicationWithWorkerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationWithWorkerServiceServer will
// result in compilation errors.
type UnsafeCommunicationWithWorkerServiceServer interface {
	mustEmbedUnimplementedCommunicationWithWorkerServiceServer()
}

func RegisterCommunicationWithWorkerServiceServer(s grpc.ServiceRegistrar, srv CommunicationWithWorkerServiceServer) {
	s.RegisterService(&CommunicationWithWorkerService_ServiceDesc, srv)
}

func _CommunicationWithWorkerService_AssignMapTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithWorkerServiceServer).AssignMapTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithWorkerService/assignMapTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithWorkerServiceServer).AssignMapTask(ctx, req.(*MapTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithWorkerService_AssignReduceTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReduceTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithWorkerServiceServer).AssignReduceTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithWorkerService/assignReduceTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithWorkerServiceServer).AssignReduceTask(ctx, req.(*ReduceTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithWorkerService_SendSelectedTournamentTreeSelection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TournamentTreeSelection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithWorkerServiceServer).SendSelectedTournamentTreeSelection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithWorkerService/sendSelectedTournamentTreeSelection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithWorkerServiceServer).SendSelectedTournamentTreeSelection(ctx, req.(*TournamentTreeSelection))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunicationWithWorkerService_RestartTournamentTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationWithWorkerServiceServer).RestartTournamentTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ds.CommunicationWithWorkerService/restartTournamentTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationWithWorkerServiceServer).RestartTournamentTree(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// CommunicationWithWorkerService_ServiceDesc is the grpc.ServiceDesc for CommunicationWithWorkerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommunicationWithWorkerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ds.CommunicationWithWorkerService",
	HandlerType: (*CommunicationWithWorkerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "assignMapTask",
			Handler:    _CommunicationWithWorkerService_AssignMapTask_Handler,
		},
		{
			MethodName: "assignReduceTask",
			Handler:    _CommunicationWithWorkerService_AssignReduceTask_Handler,
		},
		{
			MethodName: "sendSelectedTournamentTreeSelection",
			Handler:    _CommunicationWithWorkerService_SendSelectedTournamentTreeSelection_Handler,
		},
		{
			MethodName: "restartTournamentTree",
			Handler:    _CommunicationWithWorkerService_RestartTournamentTree_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Communication.proto",
}