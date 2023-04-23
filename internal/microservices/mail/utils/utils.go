package utils

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/mail/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

func ProtoByFolder(folder *models.Folder) *mail_proto.Folder {
	return &mail_proto.Folder{
		FolderID:       folder.FolderID,
		UserID:         folder.UserID,
		LocalName:      folder.LocalName,
		Name:           folder.Name,
		MessagesUnseen: int64(folder.MessagesUnseen),
		MessagesCount:  int64(folder.MessagesCount),
	}
}

func ProtoByFoldersModels(folders []models.Folder) *mail_proto.FoldersResponse {
	var protoFolders []*mail_proto.Folder
	for _, f := range folders {
		protoFolders = append(protoFolders, ProtoByFolder(&f))
	}

	return &mail_proto.FoldersResponse{Folders: protoFolders}
}

func ProtoByUserInfo(info models.UserInfo) *mail_proto.UserInfo {
	return &mail_proto.UserInfo{
		Email:     info.Email,
		FirstName: info.FirstName,
		LastName:  info.LastName,
	}
}

func ProtoByMessageInfo(info models.MessageInfo) *mail_proto.MessageInfo {
	var protoRecipients []*mail_proto.UserInfo
	for _, r := range info.Recipients {
		protoRecipients = append(protoRecipients, ProtoByUserInfo(r))
	}

	return &mail_proto.MessageInfo{
		MessageID:        info.MessageID,
		FromUser:         ProtoByUserInfo(info.FromUser),
		Recipients:       protoRecipients,
		Title:            info.Title,
		CreatedAt:        info.CreatedAt,
		Text:             info.CreatedAt,
		ReplyToMessageID: *info.ReplyToMessageID,
		ReplyTo:          ProtoByMessageInfo(*info.ReplyTo),
		Seen:             info.Seen,
		Favorite:         info.Favorite,
		Deleted:          info.Deleted,
	}
}

func ProtoMsgInfoResponseByModels(infos []models.MessageInfo) *mail_proto.MessagesInfoResponse {
	var protoMessagesInfo []*mail_proto.MessageInfo
	for _, i := range infos {
		protoMessagesInfo = append(protoMessagesInfo, ProtoByMessageInfo(i))
	}
	return &mail_proto.MessagesInfoResponse{MessagesInfo: protoMessagesInfo}
}

func ProtoMsgInfoResponseByModel(info *models.MessageInfo) *mail_proto.MessagesInfoResponse {
	var protoMessagesInfo []*mail_proto.MessageInfo
	protoMessagesInfo = append(protoMessagesInfo, ProtoByMessageInfo(*info))

	return &mail_proto.MessagesInfoResponse{MessagesInfo: protoMessagesInfo}
}

func UserInfoModelByProto(protoUserInfo *mail_proto.UserInfo) *models.UserInfo {
	return &models.UserInfo{
		FirstName: protoUserInfo.FirstName,
		LastName:  protoUserInfo.LastName,
		Email:     protoUserInfo.Email,
	}
}

func MessageInfoByProto(protoMessageInfo *mail_proto.MessageInfo) *models.MessageInfo {
	var recipients []models.UserInfo
	for _, rec := range protoMessageInfo.Recipients {
		recipients = append(recipients, *UserInfoModelByProto(rec))
	}

	var replyToInfo *models.MessageInfo = nil
	if protoMessageInfo.ReplyTo != nil {
		replyToInfo = MessageInfoByProto(protoMessageInfo.ReplyTo)
	}

	return &models.MessageInfo{
		MessageID:        protoMessageInfo.MessageID,
		FromUser:         *UserInfoModelByProto(protoMessageInfo.FromUser),
		Recipients:       recipients,
		Title:            protoMessageInfo.Title,
		CreatedAt:        protoMessageInfo.CreatedAt,
		Text:             protoMessageInfo.Text,
		ReplyToMessageID: &protoMessageInfo.ReplyToMessageID,
		ReplyTo:          replyToInfo,
		Seen:             protoMessageInfo.Seen,
		Favorite:         protoMessageInfo.Favorite,
		Deleted:          protoMessageInfo.Deleted,
	}
}

func MessagesInfoModelByProto(protoResponse *mail_proto.MessagesInfoResponse) []models.MessageInfo {
	var msgInfo []models.MessageInfo
	for _, info := range protoResponse.MessagesInfo {
		msgInfo = append(msgInfo, *MessageInfoByProto(info))
	}

	return msgInfo
}

func MessageModelByProto(protoMsg *mail_proto.Message) models.FormMessage {
	return models.FormMessage{
		Recipients:       protoMsg.Recipients,
		Title:            protoMsg.Title,
		Text:             protoMsg.Text,
		ReplyToMessageID: &protoMsg.ReplyToMessageID,
	}
}

func MessageModelByProtoSendParams(protoParams *mail_proto.SendMessageParams) models.FormMessage {
	return MessageModelByProto(protoParams.Message)
}

func ProtoSendParamsByUIDNMessage(uID uint64, form *models.FormMessage) *mail_proto.SendMessageParams {
	return &mail_proto.SendMessageParams{
		UID: uID,
		Message: &mail_proto.Message{
			Recipients:       form.Recipients,
			Title:            form.Title,
			Text:             form.Text,
			ReplyToMessageID: *form.ReplyToMessageID,
		},
	}
}

func FolderModelByProto(protoFolder *mail_proto.Folder) *models.Folder {
	return &models.Folder{
		FolderID:       protoFolder.FolderID,
		UserID:         protoFolder.UserID,
		LocalName:      protoFolder.LocalName,
		Name:           protoFolder.Name,
		MessagesUnseen: int(protoFolder.MessagesUnseen),
		MessagesCount:  int(protoFolder.MessagesCount),
	}
}

func FoldersModelByProto(protoFolders *mail_proto.FoldersResponse) []models.Folder {
	var folders []models.Folder
	for _, f := range protoFolders.Folders {
		folders = append(folders, *FolderModelByProto(f))
	}

	return folders
}
