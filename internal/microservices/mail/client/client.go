package client

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type MailClientGRPC struct {
	mailClient mail_proto.MailServiceClient
}

func NewMailClientGRPC(cc *grpc.ClientConn) mail.UseCaseI {
	return &MailClientGRPC{
		mailClient: mail_proto.NewMailServiceClient(cc),
	}
}

func (g MailClientGRPC) GetFolders(userID uint64) ([]models.Folder, error) {
	folderResp, err := g.mailClient.GetFolders(context.Background(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - GetFolders")
	}

	return utils.FoldersModelByProto(folderResp), nil
}

func (g MailClientGRPC) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	protoFolder, err := g.mailClient.GetFolderInfo(context.Background(), &mail_proto.UserFolder{
		UID:        userID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - GetFolderInfo")
	}

	return utils.FolderModelByProto(protoFolder), nil
}

func (g MailClientGRPC) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.GetFolderMessages(context.Background(), &mail_proto.UserFolder{
		UID:        userID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - GetFolderMessages")
	}

	return utils.MessagesInfoModelByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) CreateDefaultFolders(userID uint64) ([]models.Folder, error) {
	protoFolderResp, err := g.mailClient.CreateDefaultFolders(context.Background(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - CreateDefaultFolders")
	}

	return utils.FoldersModelByProto(protoFolderResp), nil
}

func (g MailClientGRPC) GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.GetMessage(context.Background(), &mail_proto.UIDMessageID{
		UID:       userID,
		MessageID: messageID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - GetMessage")
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) ValidateRecipients(recipients []string) ([]string, []string) {
	protoRc, err := g.mailClient.ValidateRecipients(context.Background(), &mail_proto.Recipients{Recipients: recipients})
	if err != nil {
		return nil, nil
	}

	return protoRc.ValidEmails, protoRc.InvalidEmails
}

func (g MailClientGRPC) SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.SendMessage(context.Background(), utils.ProtoSendParamsByUIDNMessage(userID, &message))
	if err != nil {
		return nil, errors.Wrap(err, "mail client - SendMessage")
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error {
	_, err := g.mailClient.SendFailedSendingMessage(context.Background(), &mail_proto.FailedEmailsParams{
		Recipient:     recipientEmail,
		InvalidEmails: invalidEmails,
	})
	if err != nil {
		return errors.Wrap(err, "mail client - SendFailedSendingMessage")
	}

	return nil
}

func (g MailClientGRPC) SendWelcomeMessage(recipientEmail string) error {
	_, err := g.mailClient.SendWelcomeMessage(context.Background(), &mail_proto.RecipientEmail{RecipientEmail: recipientEmail})
	if err != nil {
		return errors.Wrap(err, "mail client - SendWelcomeMessage")
	}

	return nil
}

func (g MailClientGRPC) MarkMessageAsSeen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.MarkMessageAsSeen(context.Background(), &mail_proto.UIDMessageID{
		UID:       userID,
		MessageID: messageID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mail client - MarkMessageAsSeen")
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) MarkMessageAsUnseen(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.MarkMessageAsUnseen(context.Background(), &mail_proto.UIDMessageID{
		UID:       userID,
		MessageID: messageID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "mail client - MarkMessageAsUnseen")
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}
