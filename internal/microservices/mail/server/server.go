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

func (g *AuthServerGRPC) GetFolders(ctx context.Context, protoUID *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.GetFolders(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolders")
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *AuthServerGRPC) GetFolderInfo(ctx context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.Folder, error) {
	folder, err := g.mailUC.GetFolderInfo(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolderInfo")
	}

	return utils.ProtoByFolder(folder), nil
}

func (g *AuthServerGRPC) GetFolderMessages(ctx context.Context, protoFolder *mail_proto.UserFolder) (*mail_proto.MessagesInfoResponse, error) {
	msfInfos, err := g.mailUC.GetFolderMessages(protoFolder.UID, protoFolder.FolderSlug)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetFolderMessages")
	}

	return utils.ProtoMsgInfoResponseByModels(msfInfos), nil
}

func (g *AuthServerGRPC) CreateDefaultFolders(ctx context.Context, protoUid *mail_proto.UID) (*mail_proto.FoldersResponse, error) {
	folders, err := g.mailUC.CreateDefaultFolders(protoUid.UID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - CreateDefaultFolders")
	}

	return utils.ProtoByFoldersModels(folders), nil
}

func (g *AuthServerGRPC) GetMessage(ctx context.Context, protoMId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	msfInfo, err := g.mailUC.GetMessage(protoMId.UID, protoMId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - GetMessage")
	}

	return utils.ProtoByMessageInfo(*msfInfo), nil
}

func (g *AuthServerGRPC) ValidateRecipients(ctx context.Context, protoRecipients *mail_proto.Recipients) (*mail_proto.ValidateRecipientsResponse, error) {
	valid, invalid := g.mailUC.ValidateRecipients(protoRecipients.Recipients)
	return &mail_proto.ValidateRecipientsResponse{
		ValidEmails:   valid,
		InvalidEmails: invalid,
	}, nil
}

func (g *AuthServerGRPC) SendMessage(ctx context.Context, protoParams *mail_proto.SendMessageParams) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.SendMessage(protoParams.UID, utils.MessageModelByProtoSendParams(protoParams))
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendMessage")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) SendFailedSendingMessage(ctx context.Context, protoParams *mail_proto.FailedEmailsParams) (*mail_proto.Nothing, error) {
	err := g.mailUC.SendFailedSendingMessage(protoParams.Recipient, protoParams.InvalidEmails)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendFailedSendingMessage")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) SendWelcomeMessage(ctx context.Context, protoEmail *mail_proto.RecipientEmail) (*mail_proto.Nothing, error) {
	err := g.mailUC.SendWelcomeMessage(protoEmail.RecipientEmail)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - SendWelcomeMessage")
	}

	return &mail_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) MarkMessageAsSeen(ctx context.Context, protoId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsSeen(protoId.UID, protoId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - MarkMessageAsSeen")
	}

	return utils.ProtoByMessageInfo(*info), nil
}

func (g *AuthServerGRPC) MarkMessageAsUnseen(ctx context.Context, protoId *mail_proto.UIDMessageID) (*mail_proto.MessageInfo, error) {
	info, err := g.mailUC.MarkMessageAsSeen(protoId.UID, protoId.MessageID)
	if err != nil {
		return nil, errors.Wrap(err, "mail server - MarkMessageAsUnseen")
	}

	return utils.ProtoByMessageInfo(*info), nil
}
