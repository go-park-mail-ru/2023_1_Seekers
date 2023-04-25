package client

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type UserClientGRPC struct {
	userClient user_proto.UserServiceClient
}

func NewUserClientGRPC(cc *grpc.ClientConn) user.UseCaseI {
	return &UserClientGRPC{
		userClient: user_proto.NewUserServiceClient(cc),
	}
}

func (g UserClientGRPC) Create(user *models.User) (*models.User, error) {
	protoUser, err := g.userClient.Create(context.TODO(), utils.ProtoByUserModel(user))
	if err != nil {
		return nil, errors.Wrap(err, "user client - Create")
	}

	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) Delete(ID uint64) error {
	_, err := g.userClient.Delete(context.TODO(), &user_proto.UID{UID: ID})
	if err != nil {
		return errors.Wrap(err, "user client - Delete")
	}
	return nil
}

func (g UserClientGRPC) GetByID(ID uint64) (*models.User, error) {
	protoUser, err := g.userClient.GetByID(context.TODO(), &user_proto.UID{UID: ID})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetByID")
	}
	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) GetByEmail(email string) (*models.User, error) {
	protoUser, err := g.userClient.GetByEmail(context.TODO(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetByEmail")
	}
	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) GetInfo(ID uint64) (*models.UserInfo, error) {
	protoInfo, err := g.userClient.GetInfo(context.TODO(), &user_proto.UID{UID: ID})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) GetInfoByEmail(email string) (*models.UserInfo, error) {
	protoInfo, err := g.userClient.GetInfoByEmail(context.TODO(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) EditInfo(ID uint64, info *models.UserInfo) (*models.UserInfo, error) {
	protoInfo, err := g.userClient.EditInfo(context.TODO(), &user_proto.EditInfoParams{UID: ID, EditInfo: utils.ProtoByInfoModel(info)})
	if err != nil {
		return nil, errors.Wrap(err, "user client - EditInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) EditAvatar(ID uint64, newAvatar *models.Image, isCustom bool) error {
	_, err := g.userClient.EditAvatar(context.TODO(), &user_proto.EditAvatarParams{
		UID:      ID,
		NewImage: utils.ProtoByImageModel(newAvatar),
		IsCustom: isCustom,
	})
	if err != nil {
		return errors.Wrap(err, "user client - EditAvatar")
	}
	return nil
}

func (g UserClientGRPC) GetAvatar(email string) (*models.Image, error) {
	protoImg, err := g.userClient.GetAvatar(context.TODO(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetAvatar")
	}
	return utils.ImageModelByProto(protoImg), nil
}

func (g UserClientGRPC) EditPw(ID uint64, form *models.EditPasswordRequest) error {
	_, err := g.userClient.EditPw(context.TODO(), utils.ProtoByEditPasswordModel(form, ID))
	if err != nil {
		return errors.Wrap(err, "user client - EditPw")
	}
	return nil
}
