package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type AuthServerGRPC struct {
	mail_proto.UnimplementedMailServiceServer

	grpcServer *grpc.Server
	mailUC     mail.UseCaseI
}

func NewAuthServerGRPC(g *grpc.Server, mUC mail.UseCaseI) *AuthServerGRPC {
	return &AuthServerGRPC{
		grpcServer: g,
		mailUC:     mUC,
	}
}

func (g *AuthServerGRPC) Start(url string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return err
	}
	mail_proto.RegisterMailServiceServer(g.grpcServer, g)
	return g.grpcServer.Serve(lis)
}

func (g *AuthServerGRPC) GetFolders(_ context.Context, protoUID *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.GetFolders(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolders")
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *AuthServerGRPC) GetFolderInfo(_ context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.GetFolderInfo(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolderInfo")
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *AuthServerGRPC) GetFolderMessages(_ context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.MessagesInfoResponse, error) {
	msfInfos, err := g.mailUC.GetFolderMessages(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolderMessages")
	}

	return utils.ProtoMsgInfoResponseByModels(msfInfos), nil
}

func (g *AuthServerGRPC) CreateDefaultFolders(_ context.Context, protoUid *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.CreateDefaultFolders(protoUid.UID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - CreateDefaultFolders")
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *AuthServerGRPC) GetMessage(_ context.Context, protoMId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	msfInfo, err := g.mailUC.GetMessage(protoMId.UID, protoMId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetMessage")
	}

	return utils.ProtoByMessageInfo(*msfInfo), nil
}

func (g *AuthServerGRPC) ValidateRecipients(_ context.Context, protoRecipients *mail_proto.Recipients) (*mail_proto.ValidateRecipientsResponse, error) {
	valid, invalid := g.mailUC.ValidateRecipients(protoRecipients.Recipients)
	return &mail_proto.ValidateRecipientsResponse{
		ValidEmails:   valid,
		InvalidEmails: invalid,
	}, nil
}

func (g *AuthServerGRPC) SendMessage(_ context.Context, protoParams *mail_proto.SendMessageParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SendMessage(protoParams.UID, utils.MessageModelByProtoSendParams(protoParams))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendMessage")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) SendFailedSendingMessage(_ context.Context, protoParams *mail_proto.FailedEmailsParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.SendFailedSendingMessage(protoParams.Recipient, protoParams.InvalidEmails)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendFailedSendingMessage")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) SendWelcomeMessage(_ context.Context, protoEmail *mail_proto.RecipientEmail) (*mail_proto.Nothing, error) {
	err := g.mailUC.SendWelcomeMessage(protoEmail.RecipientEmail)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendWelcomeMessage")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) MarkMessageAsSeen(_ context.Context, protoId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsSeen(protoId.UID, protoId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - MarkMessageAsSeen")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) MarkMessageAsUnseen(_ context.Context, protoId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsSeen(protoId.UID, protoId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - MarkMessageAsUnseen")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) CreateFolder(_ context.Context, protoParams *mail_proto.CreateFolderParams) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.CreateFolder(protoParams.UID, utils.FormFolderModelByProto(protoParams.FormFolder))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - CreateFolder")
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *AuthServerGRPC) DeleteFolder(_ context.Context, protoParams *mail_proto.DeleteFolderParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.DeleteFolder(protoParams.UID, protoParams.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - DeleteFolder")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) EditFolder(_ context.Context, protoParams *mail_proto.EditFolderParams) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.EditFolder(protoParams.UID, protoParams.FolderSlug, utils.FormFolderModelByProto(protoParams.FormFolder))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - EditFolder")
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *AuthServerGRPC) DeleteMessage(_ context.Context, protoParams *mail_proto.UIDMessageID) (*mail_proto.Nothing, error) {
	err := g.mailUC.DeleteMessage(protoParams.UID, protoParams.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - DeleteMessage")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) SaveDraft(_ context.Context, protoParams *mail_proto.SaveDraftParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SaveDraft(protoParams.UID, utils.MessageModelByProto(protoParams.Message))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SaveDraft")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) EditDraft(_ context.Context, protoParams *mail_proto.EditDraftParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.EditDraft(protoParams.UID, protoParams.MessageID, utils.MessageModelByProto(protoParams.Message))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - EditDraft")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) MoveMessageToFolder(_ context.Context, protoParams *mail_proto.MoveToFolderParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.MoveMessageToFolder(protoParams.UID, protoParams.MessageID, protoParams.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - MoveMessageToFolder")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) GetCustomFolders(_ context.Context, protoUid *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.GetCustomFolders(protoUid.UID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetCustomFolders")
	}

	return utils.ProtoByFoldersModels(folders), nil
}
