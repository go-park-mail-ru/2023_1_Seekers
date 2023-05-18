package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/utils"
	pkgGrpc "github.com/go-park-mail-ru/2023_1_Seekers/pkg/grpc"
	"google.golang.org/grpc"
	"net"
)

type MailServerGRPC struct {
	mail_proto.UnimplementedMailServiceServer

	grpcServer *grpc.Server
	mailUC     mail.UseCaseI
}

func NewAuthServerGRPC(g *grpc.Server, mUC mail.UseCaseI) *MailServerGRPC {
	return &MailServerGRPC{
		grpcServer: g,
		mailUC:     mUC,
	}
}

func (g *MailServerGRPC) Start(url string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return err
	}
	mail_proto.RegisterMailServiceServer(g.grpcServer, g)
	return g.grpcServer.Serve(lis)
}

func (g *MailServerGRPC) GetFolders(ctx context.Context, protoUID *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.GetFolders(protoUID.UID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *MailServerGRPC) GetFolderInfo(ctx context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.GetFolderInfo(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *MailServerGRPC) GetFolderMessages(ctx context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.MessagesInfoResponse, error) {
	msfInfos, err := g.mailUC.GetFolderMessages(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoMsgInfoResponseByModels(msfInfos), nil
}

func (g *MailServerGRPC) SearchMessages(ctx context.Context, protoParams *mail_proto.SearchMailParams) (*mail_proto.MessagesInfoResponse, error) {
	msfInfos, err := g.mailUC.SearchMessages(protoParams.UID, protoParams.FromUser, protoParams.ToUser, protoParams.Folder, protoParams.Filter)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoMsgInfoResponseByModels(msfInfos), nil
}

func (g *MailServerGRPC) SearchRecipients(ctx context.Context, protoParams *mail_proto.UID) (*mail_proto.SearchRecipientsResponse, error) {
	recipesInfo, err := g.mailUC.SearchRecipients(protoParams.UID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	var protoUsers []*mail_proto.UserInfo
	for _, r := range recipesInfo {
		protoUsers = append(protoUsers, &mail_proto.UserInfo{
			Email:     r.Email,
			FirstName: r.FirstName,
			LastName:  r.LastName,
		})
	}

	return &mail_proto.SearchRecipientsResponse{UsersInfo: protoUsers}, nil
}

func (g *MailServerGRPC) CreateDefaultFolders(ctx context.Context, protoUid *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.CreateDefaultFolders(protoUid.UID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *MailServerGRPC) GetMessage(ctx context.Context, protoMId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	msfInfo, err := g.mailUC.GetMessage(protoMId.UID, protoMId.MessageID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*msfInfo), nil
}

func (g *MailServerGRPC) GetAttachInfo(ctx context.Context, protoAttach *mail_proto.AttNUser) (*mail_proto.AttachmentInfo, error) {
	attach, err := g.mailUC.GetAttach(protoAttach.AttachID, protoAttach.UserID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoAttachByModel(attach), nil
}

func (g *MailServerGRPC) ValidateRecipients(_ context.Context, protoRecipients *mail_proto.Recipients) (*mail_proto.ValidateRecipientsResponse, error) {
	valid, invalid := g.mailUC.ValidateRecipients(protoRecipients.Recipients)
	return &mail_proto.ValidateRecipientsResponse{
		ValidEmails:   valid,
		InvalidEmails: invalid,
	}, nil
}

func (g *MailServerGRPC) SendMessage(ctx context.Context, protoParams *mail_proto.SendMessageParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SendMessage(utils.MessageModelByProtoSendParams(protoParams))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) SendFailedSendingMessage(ctx context.Context, protoParams *mail_proto.FailedEmailsParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SendFailedSendingMessage(protoParams.Recipient, protoParams.InvalidEmails)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) SendWelcomeMessage(ctx context.Context, protoEmail *mail_proto.RecipientEmail) (*mail_proto.Nothing, error) {
	err := g.mailUC.SendWelcomeMessage(protoEmail.RecipientEmail)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return &mail_proto.Nothing{}, nil
}

func (g *MailServerGRPC) MarkMessageAsSeen(ctx context.Context, protoId *mail_proto.UIDMessageIDFolderSlug) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsSeen(protoId.UID, protoId.MessageID, protoId.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) MarkMessageAsUnseen(ctx context.Context, protoId *mail_proto.UIDMessageIDFolderSlug) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsUnseen(protoId.UID, protoId.MessageID, protoId.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) CreateFolder(ctx context.Context, protoParams *mail_proto.CreateFolderParams) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.CreateFolder(protoParams.UID, utils.FormFolderModelByProto(protoParams.FormFolder))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *MailServerGRPC) DeleteFolder(ctx context.Context, protoParams *mail_proto.DeleteFolderParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.DeleteFolder(protoParams.UID, protoParams.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return &mail_proto.Nothing{}, nil
}

func (g *MailServerGRPC) EditFolder(ctx context.Context, protoParams *mail_proto.EditFolderParams) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.EditFolder(protoParams.UID, protoParams.FolderSlug, utils.FormFolderModelByProto(protoParams.FormFolder))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *MailServerGRPC) DeleteMessage(ctx context.Context, protoParams *mail_proto.UIDMessageIDFolderSlug) (*mail_proto.Nothing, error) {
	err := g.mailUC.DeleteMessage(protoParams.UID, protoParams.MessageID, protoParams.FolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return &mail_proto.Nothing{}, nil
}

func (g *MailServerGRPC) SaveDraft(ctx context.Context, protoParams *mail_proto.SaveDraftParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SaveDraft(protoParams.UID, utils.MessageModelByProto(protoParams.Message))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) EditDraft(ctx context.Context, protoParams *mail_proto.EditDraftParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.EditDraft(protoParams.UID, protoParams.MessageID, utils.MessageModelByProto(protoParams.Message))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *MailServerGRPC) MoveMessageToFolder(ctx context.Context, protoParams *mail_proto.MoveToFolderParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.MoveMessageToFolder(protoParams.UID, protoParams.MessageID, protoParams.FromFolderSlug, protoParams.ToFolderSlug)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return &mail_proto.Nothing{}, nil
}

func (g *MailServerGRPC) GetCustomFolders(ctx context.Context, protoUid *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.GetCustomFolders(protoUid.UID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *MailServerGRPC) GetAttach(ctx context.Context, protoAttach *mail_proto.AttNUser) (*mail_proto.AttachmentInfo, error) {
	attach, err := g.mailUC.GetAttach(protoAttach.AttachID, protoAttach.UserID)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoAttachByModel(attach), nil
}
