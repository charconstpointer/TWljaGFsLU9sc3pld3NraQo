// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package fetcher

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// FetcherServiceClient is the client API for FetcherService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FetcherServiceClient interface {
	GetMeasures(ctx context.Context, in *GetMeasuresRequest, opts ...grpc.CallOption) (*GetMeasuresResponse, error)
	AddProbe(ctx context.Context, in *AddProbeRequest, opts ...grpc.CallOption) (*AddProbeResponse, error)
	ListenForChanges(ctx context.Context, in *ListenForChangesRequest, opts ...grpc.CallOption) (FetcherService_ListenForChangesClient, error)
}

type fetcherServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFetcherServiceClient(cc grpc.ClientConnInterface) FetcherServiceClient {
	return &fetcherServiceClient{cc}
}

func (c *fetcherServiceClient) GetMeasures(ctx context.Context, in *GetMeasuresRequest, opts ...grpc.CallOption) (*GetMeasuresResponse, error) {
	out := new(GetMeasuresResponse)
	err := c.cc.Invoke(ctx, "/fetcher.FetcherService/GetMeasures", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fetcherServiceClient) AddProbe(ctx context.Context, in *AddProbeRequest, opts ...grpc.CallOption) (*AddProbeResponse, error) {
	out := new(AddProbeResponse)
	err := c.cc.Invoke(ctx, "/fetcher.FetcherService/AddProbe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fetcherServiceClient) ListenForChanges(ctx context.Context, in *ListenForChangesRequest, opts ...grpc.CallOption) (FetcherService_ListenForChangesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_FetcherService_serviceDesc.Streams[0], "/fetcher.FetcherService/ListenForChanges", opts...)
	if err != nil {
		return nil, err
	}
	x := &fetcherServiceListenForChangesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FetcherService_ListenForChangesClient interface {
	Recv() (*ListenForChangesResponse, error)
	grpc.ClientStream
}

type fetcherServiceListenForChangesClient struct {
	grpc.ClientStream
}

func (x *fetcherServiceListenForChangesClient) Recv() (*ListenForChangesResponse, error) {
	m := new(ListenForChangesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FetcherServiceServer is the server API for FetcherService service.
// All implementations must embed UnimplementedFetcherServiceServer
// for forward compatibility
type FetcherServiceServer interface {
	GetMeasures(context.Context, *GetMeasuresRequest) (*GetMeasuresResponse, error)
	AddProbe(context.Context, *AddProbeRequest) (*AddProbeResponse, error)
	ListenForChanges(*ListenForChangesRequest, FetcherService_ListenForChangesServer) error
	mustEmbedUnimplementedFetcherServiceServer()
}

// UnimplementedFetcherServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFetcherServiceServer struct {
}

func (*UnimplementedFetcherServiceServer) GetMeasures(context.Context, *GetMeasuresRequest) (*GetMeasuresResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeasures not implemented")
}
func (*UnimplementedFetcherServiceServer) AddProbe(context.Context, *AddProbeRequest) (*AddProbeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProbe not implemented")
}
func (*UnimplementedFetcherServiceServer) ListenForChanges(*ListenForChangesRequest, FetcherService_ListenForChangesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListenForChanges not implemented")
}
func (*UnimplementedFetcherServiceServer) mustEmbedUnimplementedFetcherServiceServer() {}

func RegisterFetcherServiceServer(s *grpc.Server, srv FetcherServiceServer) {
	s.RegisterService(&_FetcherService_serviceDesc, srv)
}

func _FetcherService_GetMeasures_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeasuresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FetcherServiceServer).GetMeasures(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fetcher.FetcherService/GetMeasures",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FetcherServiceServer).GetMeasures(ctx, req.(*GetMeasuresRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FetcherService_AddProbe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProbeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FetcherServiceServer).AddProbe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fetcher.FetcherService/AddProbe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FetcherServiceServer).AddProbe(ctx, req.(*AddProbeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FetcherService_ListenForChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenForChangesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FetcherServiceServer).ListenForChanges(m, &fetcherServiceListenForChangesServer{stream})
}

type FetcherService_ListenForChangesServer interface {
	Send(*ListenForChangesResponse) error
	grpc.ServerStream
}

type fetcherServiceListenForChangesServer struct {
	grpc.ServerStream
}

func (x *fetcherServiceListenForChangesServer) Send(m *ListenForChangesResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _FetcherService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fetcher.FetcherService",
	HandlerType: (*FetcherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMeasures",
			Handler:    _FetcherService_GetMeasures_Handler,
		},
		{
			MethodName: "AddProbe",
			Handler:    _FetcherService_AddProbe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenForChanges",
			Handler:       _FetcherService_ListenForChanges_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/fetcher/pb/fetcher.proto",
}
