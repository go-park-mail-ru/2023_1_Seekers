package client

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	pkgGrpc "github.com/go-park-mail-ru/2023_1_Seekers/pkg/grpc"
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
	folderResp, err := g.mailClient.GetFolders(context.TODO(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - GetFolders"))
	}

	return utils.FoldersModelByProto(folderResp), nil
}

func (g MailClientGRPC) GetFolderInfo(userID uint64, folderSlug string) (*models.Folder, error) {
	protoFolder, err := g.mailClient.GetFolderInfo(context.TODO(), &mail_proto.UserFolder{
		UID:        userID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - GetFolderInfo"))
	}

	return utils.FolderModelByProto(protoFolder), nil
}

func (g MailClientGRPC) GetFolderMessages(userID uint64, folderSlug string) ([]models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.GetFolderMessages(context.TODO(), &mail_proto.UserFolder{
		UID:        userID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - GetFolderMessages"))
	}

	return utils.MessagesInfoModelByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) SearchMessages(userId uint64, fromUser, toUser, folder, filter string) ([]models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.SearchMessages(context.TODO(), &mail_proto.SearchMailParams{
		UID:      userId,
		FromUser: fromUser,
		ToUser:   toUser,
		Folder:   folder,
		Filter:   filter,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - SearchMessages"))
	}

	return utils.MessagesInfoModelByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) SearchRecipients(userID uint64) ([]models.UserInfo, error) {
	protoUsers, err := g.mailClient.SearchRecipients(context.TODO(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - SearchRecipients"))
	}

	var recipesInfo []models.UserInfo
	for _, r := range protoUsers.UsersInfo {
		recipesInfo = append(recipesInfo, models.UserInfo{
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Email:     r.Email,
		})
	}

	return recipesInfo, nil
}

func (g MailClientGRPC) CreateDefaultFolders(userID uint64) ([]models.Folder, error) {
	protoFolderResp, err := g.mailClient.CreateDefaultFolders(context.TODO(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - CreateDefaultFolders"))
	}

	return utils.FoldersModelByProto(protoFolderResp), nil
}

func (g MailClientGRPC) GetMessage(userID uint64, messageID uint64) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.GetMessage(context.TODO(), &mail_proto.UIDMessageID{
		UID:       userID,
		MessageID: messageID,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - GetMessage"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) ValidateRecipients(recipients []string) ([]string, []string) {
	protoRc, err := g.mailClient.ValidateRecipients(context.TODO(), &mail_proto.Recipients{Recipients: recipients})
	if err != nil {
		return nil, nil
	}

	return protoRc.ValidEmails, protoRc.InvalidEmails
}

func (g MailClientGRPC) SendMessage(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.SendMessage(context.TODO(), utils.ProtoSendParamsByUIDNMessage(userID, &message))
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - SendMessage"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) SendFailedSendingMessage(recipientEmail string, invalidEmails []string) error {
	_, err := g.mailClient.SendFailedSendingMessage(context.TODO(), &mail_proto.FailedEmailsParams{
		Recipient:     recipientEmail,
		InvalidEmails: invalidEmails,
	})
	if err != nil {
		return pkgGrpc.CauseError(errors.Wrap(err, "mail client - SendFailedSendingMessage"))
	}

	return nil
}

func (g MailClientGRPC) SendWelcomeMessage(recipientEmail string) error {
	_, err := g.mailClient.SendWelcomeMessage(context.TODO(), &mail_proto.RecipientEmail{RecipientEmail: recipientEmail})
	if err != nil {
		return pkgGrpc.CauseError(errors.Wrap(err, "mail client - SendWelcomeMessage"))
	}

	return nil
}

func (g MailClientGRPC) MarkMessageAsSeen(userID uint64, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.MarkMessageAsSeen(context.TODO(), &mail_proto.UIDMessageIDFolderSlug{
		UID:        userID,
		MessageID:  messageID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - MarkMessageAsSeen"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) MarkMessageAsUnseen(userID uint64, messageID uint64, folderSlug string) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.MarkMessageAsUnseen(context.TODO(), &mail_proto.UIDMessageIDFolderSlug{
		UID:        userID,
		MessageID:  messageID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - MarkMessageAsUnseen"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) CreateFolder(userID uint64, form models.FormFolder) (*models.Folder, error) {
	protoFolder, err := g.mailClient.CreateFolder(context.TODO(), &mail_proto.CreateFolderParams{
		UID:        userID,
		FormFolder: &mail_proto.FormFolder{Name: form.Name},
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - CreateFolder"))
	}

	return utils.FolderModelByProto(protoFolder), nil
}

func (g MailClientGRPC) DeleteFolder(userID uint64, folderSlug string) error {
	_, err := g.mailClient.DeleteFolder(context.TODO(), &mail_proto.DeleteFolderParams{
		UID:        userID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return pkgGrpc.CauseError(errors.Wrap(err, "mail client - DeleteFolder"))
	}

	return nil
}

func (g MailClientGRPC) EditFolder(userID uint64, folderSlug string, form models.FormFolder) (*models.Folder, error) {
	protoFolder, err := g.mailClient.EditFolder(context.TODO(), &mail_proto.EditFolderParams{
		UID:        userID,
		FolderSlug: folderSlug,
		FormFolder: &mail_proto.FormFolder{Name: form.Name},
	})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - EditFolder"))
	}

	return utils.FolderModelByProto(protoFolder), nil
}

func (g MailClientGRPC) DeleteMessage(userID uint64, messageID uint64, folderSlug string) error {
	_, err := g.mailClient.DeleteMessage(context.TODO(), &mail_proto.UIDMessageIDFolderSlug{
		UID:        userID,
		MessageID:  messageID,
		FolderSlug: folderSlug,
	})
	if err != nil {
		return pkgGrpc.CauseError(errors.Wrap(err, "mail client - DeleteMessage"))
	}

	return nil
}

func (g MailClientGRPC) SaveDraft(userID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.SaveDraft(context.TODO(), utils.ProtoSaveDraftParamsByModels(userID, &message))
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - SaveDraft"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) EditDraft(userID, messageID uint64, message models.FormMessage) (*models.MessageInfo, error) {
	protoMsgInfo, err := g.mailClient.EditDraft(context.TODO(), utils.ProtoEditDraftParamsByModels(userID, messageID, &message))
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - EditDraft"))
	}

	return utils.MessageInfoByProto(protoMsgInfo), nil
}

func (g MailClientGRPC) MoveMessageToFolder(userID, messageID uint64, fromFolderSlug string, toFolderSlug string) error {
	_, err := g.mailClient.MoveMessageToFolder(context.TODO(), &mail_proto.MoveToFolderParams{
		UID:            userID,
		MessageID:      messageID,
		FromFolderSlug: fromFolderSlug,
		ToFolderSlug:   toFolderSlug,
	})
	if err != nil {
		return pkgGrpc.CauseError(errors.Wrap(err, "mail client - MoveMessageToFolder"))
	}

	return nil
}

func (g MailClientGRPC) GetCustomFolders(userID uint64) ([]models.Folder, error) {
	protoFolderResp, err := g.mailClient.GetCustomFolders(context.TODO(), &mail_proto.UID{UID: userID})
	if err != nil {
		return nil, pkgGrpc.CauseError(errors.Wrap(err, "mail client - GetCustomFolders"))
	}

	return utils.FoldersModelByProto(protoFolderResp), nil
}
