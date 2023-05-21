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

func ModelAttachByProto(attach *mail_proto.AttachmentInfo) *models.AttachmentInfo {
	return &models.AttachmentInfo{
		AttachID:  attach.AttachID,
		FileName:  attach.FileName,
		FileData:  attach.FileData,
		Type:      attach.Type,
		S3FName:   attach.S3FName,
		SizeCount: attach.SizeCount,
		SizeStr:   attach.SizeStr,
	}
}

func ProtoAttachByModel(info *models.AttachmentInfo) *mail_proto.AttachmentInfo {
	return &mail_proto.AttachmentInfo{
		AttachID:  info.AttachID,
		FileName:  info.FileName,
		FileData:  info.FileData,
		Type:      info.Type,
		S3FName:   info.S3FName,
		SizeCount: info.SizeCount,
		SizeStr:   info.SizeStr,
	}
}

func ProtoByMessageInfo(info models.MessageInfo) *mail_proto.MessageInfo {
	var protoRecipients []*mail_proto.UserInfo
	for _, r := range info.Recipients {
		protoRecipients = append(protoRecipients, ProtoByUserInfo(r))
	}

	var replyToInfo *mail_proto.MessageInfo = nil
	if info.ReplyTo != nil {
		replyToInfo = ProtoByMessageInfo(*info.ReplyTo)
	}

	var replyMessageID *mail_proto.UID = nil
	if info.ReplyToMessageID != nil {
		replyMessageID = &mail_proto.UID{UID: *info.ReplyToMessageID}
	}

	protoAttaches := make([]*mail_proto.AttachmentInfo, len(info.Attachments))
	for i, a := range info.Attachments {
		protoAttaches[i] = ProtoAttachByModel(&a)
	}

	return &mail_proto.MessageInfo{
		MessageID:        info.MessageID,
		FromUser:         ProtoByUserInfo(info.FromUser),
		Recipients:       protoRecipients,
		Title:            info.Title,
		CreatedAt:        info.CreatedAt,
		Text:             info.Text,
		ReplyToMessageID: replyMessageID,
		ReplyTo:          replyToInfo,
		Seen:             info.Seen,
		Favorite:         info.Favorite,
		Deleted:          info.Deleted,
		Attachments:      protoAttaches,
		AttachmentsSize:  info.AttachmentsSize,
		Preview:          info.Preview,
	}
}

func ProtoMsgInfoResponseByModels(infos []models.MessageInfo) *mail_proto.MessagesInfoResponse {
	var protoMessagesInfo []*mail_proto.MessageInfo
	for _, i := range infos {
		protoMessagesInfo = append(protoMessagesInfo, ProtoByMessageInfo(i))
	}
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

	var replyMessageID *uint64 = nil
	if protoMessageInfo.ReplyToMessageID != nil {
		replyMessageID = &protoMessageInfo.ReplyToMessageID.UID
	}

	attaches := make([]models.AttachmentInfo, len(protoMessageInfo.Attachments))
	for i, a := range protoMessageInfo.Attachments {
		attaches[i] = *ModelAttachByProto(a)
	}

	return &models.MessageInfo{
		MessageID:        protoMessageInfo.MessageID,
		FromUser:         *UserInfoModelByProto(protoMessageInfo.FromUser),
		Recipients:       recipients,
		Title:            protoMessageInfo.Title,
		CreatedAt:        protoMessageInfo.CreatedAt,
		Text:             protoMessageInfo.Text,
		ReplyToMessageID: replyMessageID,
		ReplyTo:          replyToInfo,
		Seen:             protoMessageInfo.Seen,
		Favorite:         protoMessageInfo.Favorite,
		Deleted:          protoMessageInfo.Deleted,
		Attachments:      attaches,
		AttachmentsSize:  protoMessageInfo.AttachmentsSize,
		Preview:          protoMessageInfo.Preview,
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
	var replyMessageID *uint64 = nil
	if protoMsg.ReplyToMessageID != nil {
		replyMessageID = &protoMsg.ReplyToMessageID.UID
	}

	attaches := make([]models.Attachment, len(protoMsg.Attachments))
	for i, a := range protoMsg.Attachments {
		attaches[i] = models.Attachment{
			FileName: a.FileName,
			FileData: a.FileData,
		}
	}

	return models.FormMessage{
		FromUser:         protoMsg.FromUser,
		Recipients:       protoMsg.Recipients,
		Title:            protoMsg.Title,
		Text:             protoMsg.Text,
		ReplyToMessageID: replyMessageID,
		Attachments:      attaches,
	}
}

func MessageModelByProtoSendParams(protoParams *mail_proto.SendMessageParams) (uint64, models.FormMessage) {
	return protoParams.UID, MessageModelByProto(protoParams.Message)
}

func ProtoSendParamsByUIDNMessage(userId uint64, form *models.FormMessage) *mail_proto.SendMessageParams {
	var replyMessageID *mail_proto.UID = nil
	if form.ReplyToMessageID != nil {
		replyMessageID = &mail_proto.UID{UID: *form.ReplyToMessageID}
	}

	attaches := make([]*mail_proto.Attachment, len(form.Attachments))
	for i, a := range form.Attachments {
		attaches[i] = &mail_proto.Attachment{
			FileName: a.FileName,
			FileData: a.FileData,
		}
	}

	return &mail_proto.SendMessageParams{
		UID: userId,
		Message: &mail_proto.Message{
			FromUser:         form.FromUser,
			Recipients:       form.Recipients,
			Title:            form.Title,
			Text:             form.Text,
			ReplyToMessageID: replyMessageID,
			Attachments:      attaches,
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

func ProtoSaveDraftParamsByModels(uID uint64, form *models.FormMessage) *mail_proto.SaveDraftParams {
	var replyMessageID *mail_proto.UID = nil
	if form.ReplyToMessageID != nil {
		replyMessageID = &mail_proto.UID{UID: *form.ReplyToMessageID}
	}

	attaches := make([]*mail_proto.Attachment, len(form.Attachments))
	for i, a := range form.Attachments {
		attaches[i] = &mail_proto.Attachment{
			FileName: a.FileName,
			FileData: a.FileData,
		}
	}

	return &mail_proto.SaveDraftParams{
		UID: uID,
		Message: &mail_proto.Message{
			Recipients:       form.Recipients,
			Title:            form.Title,
			Text:             form.Text,
			ReplyToMessageID: replyMessageID,
			Attachments:      attaches,
		},
	}
}

func ProtoEditDraftParamsByModels(uID, messageId uint64, form *models.FormMessage) *mail_proto.EditDraftParams {
	var replyMessageID *mail_proto.UID = nil
	if form.ReplyToMessageID != nil {
		replyMessageID = &mail_proto.UID{UID: *form.ReplyToMessageID}
	}

	return &mail_proto.EditDraftParams{
		UID:       uID,
		MessageID: messageId,
		Message: &mail_proto.Message{
			Recipients:       form.Recipients,
			Title:            form.Title,
			Text:             form.Text,
			ReplyToMessageID: replyMessageID,
		},
	}
}

func FormFolderModelByProto(protoForm *mail_proto.FormFolder) models.FormFolder {
	return models.FormFolder{Name: protoForm.Name}
}
