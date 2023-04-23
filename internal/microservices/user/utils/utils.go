package utils

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

func UserModelByProto(proto *user_proto.User) *models.User {
	return &models.User{
		UserID:    proto.UID,
		Email:     proto.Email,
		Password:  string(proto.Password),
		FirstName: proto.FirstName,
		LastName:  proto.LastName,
		Avatar:    proto.Avatar,
	}
}

func ProtoByUserModel(user *models.User) *user_proto.User {
	return &user_proto.User{
		UID:       user.UserID,
		Email:     user.Email,
		Password:  []byte(user.Password),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
}

func ProtoByInfoModel(info *models.UserInfo) *user_proto.UserInfo {
	return &user_proto.UserInfo{
		UID:       info.UserID,
		Email:     info.Email,
		FirstName: info.FirstName,
		LastName:  info.LastName,
	}
}

func InfoModelByProto(protoInfo *user_proto.UserInfo) *models.UserInfo {
	return &models.UserInfo{
		UserID:    protoInfo.UID,
		FirstName: protoInfo.FirstName,
		LastName:  protoInfo.LastName,
		Email:     protoInfo.Email,
	}
}

func ImageModelByProto(protoImg *user_proto.Image) *models.Image {
	return &models.Image{
		Name: protoImg.Name,
		Data: protoImg.Data,
	}
}

func ProtoByImageModel(img *models.Image) *user_proto.Image {
	return &user_proto.Image{
		Name: img.Name,
		Data: img.Data,
	}
}

func EditPasswordModelByProto(protoForm *user_proto.EditPasswordParams) *models.EditPasswordRequest {
	return &models.EditPasswordRequest{
		PasswordOld: protoForm.PasswordOld,
		Password:    protoForm.Password,
		RepeatPw:    protoForm.RepeatPw,
	}
}

func ProtoByEditPasswordModel(form *models.EditPasswordRequest, ID uint64) *user_proto.EditPasswordParams {
	return &user_proto.EditPasswordParams{
		UID:         ID,
		PasswordOld: form.PasswordOld,
		Password:    form.Password,
		RepeatPw:    form.RepeatPw,
	}
}
