// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mail_proto

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

// MailServiceClient is the client API for MailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MailServiceClient interface {
	GetFolders(ctx context.Context, in *UID, opts ...grpc.CallOption) (*FoldersResponse, error)
	GetFolderInfo(ctx context.Context, in *UserFolder, opts ...grpc.CallOption) (*Folder, error)
	GetFolderMessages(ctx context.Context, in *UserFolder, opts ...grpc.CallOption) (*MessagesInfoResponse, error)
	CreateDefaultFolders(ctx context.Context, in *UID, opts ...grpc.CallOption) (*FoldersResponse, error)
	GetMessage(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error)
	ValidateRecipients(ctx context.Context, in *Recipients, opts ...grpc.CallOption) (*ValidateRecipientsResponse, error)
	SendMessage(ctx context.Context, in *SendMessageParams, opts ...grpc.CallOption) (*MessageInfo, error)
	SendFailedSendingMessage(ctx context.Context, in *FailedEmailsParams, opts ...grpc.CallOption) (*Nothing, error)
	SendWelcomeMessage(ctx context.Context, in *RecipientEmail, opts ...grpc.CallOption) (*Nothing, error)
	MarkMessageAsSeen(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error)
	MarkMessageAsUnseen(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error)
	CreateFolder(ctx context.Context, in *CreateFolderParams, opts ...grpc.CallOption) (*Folder, error)
	DeleteFolder(ctx context.Context, in *DeleteFolderParams, opts ...grpc.CallOption) (*Nothing, error)
	EditFolder(ctx context.Context, in *EditFolderParams, opts ...grpc.CallOption) (*Folder, error)
	DeleteMessage(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*Nothing, error)
	SaveDraft(ctx context.Context, in *SaveDraftParams, opts ...grpc.CallOption) (*MessageInfo, error)
	EditDraft(ctx context.Context, in *EditDraftParams, opts ...grpc.CallOption) (*MessageInfo, error)
	MoveMessageToFolder(ctx context.Context, in *MoveToFolderParams, opts ...grpc.CallOption) (*Nothing, error)
}

type mailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMailServiceClient(cc grpc.ClientConnInterface) MailServiceClient {
	return &mailServiceClient{cc}
}

func (c *mailServiceClient) GetFolders(ctx context.Context, in *UID, opts ...grpc.CallOption) (*FoldersResponse, error) {
	out := new(FoldersResponse)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/GetFolders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) GetFolderInfo(ctx context.Context, in *UserFolder, opts ...grpc.CallOption) (*Folder, error) {
	out := new(Folder)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/GetFolderInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) GetFolderMessages(ctx context.Context, in *UserFolder, opts ...grpc.CallOption) (*MessagesInfoResponse, error) {
	out := new(MessagesInfoResponse)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/GetFolderMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) CreateDefaultFolders(ctx context.Context, in *UID, opts ...grpc.CallOption) (*FoldersResponse, error) {
	out := new(FoldersResponse)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/CreateDefaultFolders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) GetMessage(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/GetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) ValidateRecipients(ctx context.Context, in *Recipients, opts ...grpc.CallOption) (*ValidateRecipientsResponse, error) {
	out := new(ValidateRecipientsResponse)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/ValidateRecipients", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) SendMessage(ctx context.Context, in *SendMessageParams, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) SendFailedSendingMessage(ctx context.Context, in *FailedEmailsParams, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/SendFailedSendingMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) SendWelcomeMessage(ctx context.Context, in *RecipientEmail, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/SendWelcomeMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) MarkMessageAsSeen(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/MarkMessageAsSeen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) MarkMessageAsUnseen(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/MarkMessageAsUnseen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) CreateFolder(ctx context.Context, in *CreateFolderParams, opts ...grpc.CallOption) (*Folder, error) {
	out := new(Folder)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/CreateFolder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) DeleteFolder(ctx context.Context, in *DeleteFolderParams, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/DeleteFolder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) EditFolder(ctx context.Context, in *EditFolderParams, opts ...grpc.CallOption) (*Folder, error) {
	out := new(Folder)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/EditFolder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) DeleteMessage(ctx context.Context, in *UIDMessageID, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/DeleteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) SaveDraft(ctx context.Context, in *SaveDraftParams, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/SaveDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) EditDraft(ctx context.Context, in *EditDraftParams, opts ...grpc.CallOption) (*MessageInfo, error) {
	out := new(MessageInfo)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/EditDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailServiceClient) MoveMessageToFolder(ctx context.Context, in *MoveToFolderParams, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/mail_proto.MailService/MoveMessageToFolder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MailServiceServer is the server API for MailService service.
// All implementations must embed UnimplementedMailServiceServer
// for forward compatibility
type MailServiceServer interface {
	GetFolders(context.Context, *UID) (*FoldersResponse, error)
	GetFolderInfo(context.Context, *UserFolder) (*Folder, error)
	GetFolderMessages(context.Context, *UserFolder) (*MessagesInfoResponse, error)
	CreateDefaultFolders(context.Context, *UID) (*FoldersResponse, error)
	GetMessage(context.Context, *UIDMessageID) (*MessageInfo, error)
	ValidateRecipients(context.Context, *Recipients) (*ValidateRecipientsResponse, error)
	SendMessage(context.Context, *SendMessageParams) (*MessageInfo, error)
	SendFailedSendingMessage(context.Context, *FailedEmailsParams) (*Nothing, error)
	SendWelcomeMessage(context.Context, *RecipientEmail) (*Nothing, error)
	MarkMessageAsSeen(context.Context, *UIDMessageID) (*MessageInfo, error)
	MarkMessageAsUnseen(context.Context, *UIDMessageID) (*MessageInfo, error)
	CreateFolder(context.Context, *CreateFolderParams) (*Folder, error)
	DeleteFolder(context.Context, *DeleteFolderParams) (*Nothing, error)
	EditFolder(context.Context, *EditFolderParams) (*Folder, error)
	DeleteMessage(context.Context, *UIDMessageID) (*Nothing, error)
	SaveDraft(context.Context, *SaveDraftParams) (*MessageInfo, error)
	EditDraft(context.Context, *EditDraftParams) (*MessageInfo, error)
	MoveMessageToFolder(context.Context, *MoveToFolderParams) (*Nothing, error)
	mustEmbedUnimplementedMailServiceServer()
}

// UnimplementedMailServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMailServiceServer struct {
}

func (UnimplementedMailServiceServer) GetFolders(context.Context, *UID) (*FoldersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFolders not implemented")
}
func (UnimplementedMailServiceServer) GetFolderInfo(context.Context, *UserFolder) (*Folder, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFolderInfo not implemented")
}
func (UnimplementedMailServiceServer) GetFolderMessages(context.Context, *UserFolder) (*MessagesInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFolderMessages not implemented")
}
func (UnimplementedMailServiceServer) CreateDefaultFolders(context.Context, *UID) (*FoldersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDefaultFolders not implemented")
}
func (UnimplementedMailServiceServer) GetMessage(context.Context, *UIDMessageID) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedMailServiceServer) ValidateRecipients(context.Context, *Recipients) (*ValidateRecipientsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateRecipients not implemented")
}
func (UnimplementedMailServiceServer) SendMessage(context.Context, *SendMessageParams) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMailServiceServer) SendFailedSendingMessage(context.Context, *FailedEmailsParams) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFailedSendingMessage not implemented")
}
func (UnimplementedMailServiceServer) SendWelcomeMessage(context.Context, *RecipientEmail) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendWelcomeMessage not implemented")
}
func (UnimplementedMailServiceServer) MarkMessageAsSeen(context.Context, *UIDMessageID) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkMessageAsSeen not implemented")
}
func (UnimplementedMailServiceServer) MarkMessageAsUnseen(context.Context, *UIDMessageID) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkMessageAsUnseen not implemented")
}
func (UnimplementedMailServiceServer) CreateFolder(context.Context, *CreateFolderParams) (*Folder, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFolder not implemented")
}
func (UnimplementedMailServiceServer) DeleteFolder(context.Context, *DeleteFolderParams) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFolder not implemented")
}
func (UnimplementedMailServiceServer) EditFolder(context.Context, *EditFolderParams) (*Folder, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditFolder not implemented")
}
func (UnimplementedMailServiceServer) DeleteMessage(context.Context, *UIDMessageID) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedMailServiceServer) SaveDraft(context.Context, *SaveDraftParams) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveDraft not implemented")
}
func (UnimplementedMailServiceServer) EditDraft(context.Context, *EditDraftParams) (*MessageInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditDraft not implemented")
}
func (UnimplementedMailServiceServer) MoveMessageToFolder(context.Context, *MoveToFolderParams) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveMessageToFolder not implemented")
}
func (UnimplementedMailServiceServer) mustEmbedUnimplementedMailServiceServer() {}

// UnsafeMailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MailServiceServer will
// result in compilation errors.
type UnsafeMailServiceServer interface {
	mustEmbedUnimplementedMailServiceServer()
}

func RegisterMailServiceServer(s grpc.ServiceRegistrar, srv MailServiceServer) {
	s.RegisterService(&MailService_ServiceDesc, srv)
}

func _MailService_GetFolders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).GetFolders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/GetFolders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).GetFolders(ctx, req.(*UID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_GetFolderInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFolder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).GetFolderInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/GetFolderInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).GetFolderInfo(ctx, req.(*UserFolder))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_GetFolderMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFolder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).GetFolderMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/GetFolderMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).GetFolderMessages(ctx, req.(*UserFolder))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_CreateDefaultFolders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).CreateDefaultFolders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/CreateDefaultFolders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).CreateDefaultFolders(ctx, req.(*UID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UIDMessageID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/GetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).GetMessage(ctx, req.(*UIDMessageID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_ValidateRecipients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Recipients)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).ValidateRecipients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/ValidateRecipients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).ValidateRecipients(ctx, req.(*Recipients))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).SendMessage(ctx, req.(*SendMessageParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_SendFailedSendingMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FailedEmailsParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).SendFailedSendingMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/SendFailedSendingMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).SendFailedSendingMessage(ctx, req.(*FailedEmailsParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_SendWelcomeMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientEmail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).SendWelcomeMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/SendWelcomeMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).SendWelcomeMessage(ctx, req.(*RecipientEmail))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_MarkMessageAsSeen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UIDMessageID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).MarkMessageAsSeen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/MarkMessageAsSeen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).MarkMessageAsSeen(ctx, req.(*UIDMessageID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_MarkMessageAsUnseen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UIDMessageID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).MarkMessageAsUnseen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/MarkMessageAsUnseen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).MarkMessageAsUnseen(ctx, req.(*UIDMessageID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_CreateFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFolderParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).CreateFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/CreateFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).CreateFolder(ctx, req.(*CreateFolderParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_DeleteFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFolderParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).DeleteFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/DeleteFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).DeleteFolder(ctx, req.(*DeleteFolderParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_EditFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditFolderParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).EditFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/EditFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).EditFolder(ctx, req.(*EditFolderParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UIDMessageID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).DeleteMessage(ctx, req.(*UIDMessageID))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_SaveDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveDraftParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).SaveDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/SaveDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).SaveDraft(ctx, req.(*SaveDraftParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_EditDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditDraftParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).EditDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/EditDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).EditDraft(ctx, req.(*EditDraftParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailService_MoveMessageToFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveToFolderParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailServiceServer).MoveMessageToFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mail_proto.MailService/MoveMessageToFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailServiceServer).MoveMessageToFolder(ctx, req.(*MoveToFolderParams))
	}
	return interceptor(ctx, in, info, handler)
}

// MailService_ServiceDesc is the grpc.ServiceDesc for MailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mail_proto.MailService",
	HandlerType: (*MailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFolders",
			Handler:    _MailService_GetFolders_Handler,
		},
		{
			MethodName: "GetFolderInfo",
			Handler:    _MailService_GetFolderInfo_Handler,
		},
		{
			MethodName: "GetFolderMessages",
			Handler:    _MailService_GetFolderMessages_Handler,
		},
		{
			MethodName: "CreateDefaultFolders",
			Handler:    _MailService_CreateDefaultFolders_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _MailService_GetMessage_Handler,
		},
		{
			MethodName: "ValidateRecipients",
			Handler:    _MailService_ValidateRecipients_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _MailService_SendMessage_Handler,
		},
		{
			MethodName: "SendFailedSendingMessage",
			Handler:    _MailService_SendFailedSendingMessage_Handler,
		},
		{
			MethodName: "SendWelcomeMessage",
			Handler:    _MailService_SendWelcomeMessage_Handler,
		},
		{
			MethodName: "MarkMessageAsSeen",
			Handler:    _MailService_MarkMessageAsSeen_Handler,
		},
		{
			MethodName: "MarkMessageAsUnseen",
			Handler:    _MailService_MarkMessageAsUnseen_Handler,
		},
		{
			MethodName: "CreateFolder",
			Handler:    _MailService_CreateFolder_Handler,
		},
		{
			MethodName: "DeleteFolder",
			Handler:    _MailService_DeleteFolder_Handler,
		},
		{
			MethodName: "EditFolder",
			Handler:    _MailService_EditFolder_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _MailService_DeleteMessage_Handler,
		},
		{
			MethodName: "SaveDraft",
			Handler:    _MailService_SaveDraft_Handler,
		},
		{
			MethodName: "EditDraft",
			Handler:    _MailService_EditDraft_Handler,
		},
		{
			MethodName: "MoveMessageToFolder",
			Handler:    _MailService_MoveMessageToFolder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mail.proto",
}
