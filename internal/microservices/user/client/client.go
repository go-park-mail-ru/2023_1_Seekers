package client

import (
	"context"
	"fmt"
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
	fmt.Println(utils.ProtoByUserModel(user))
	protoUser, err := g.userClient.Create(context.Background(), utils.ProtoByUserModel(user))
	fmt.Println("_____")
	fmt.Println(protoUser)
	if err != nil {
		return nil, errors.Wrap(err, "user client - Create")
	}
	fmt.Println(utils.UserModelByProto(protoUser))
	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) Delete(ID uint64) error {
	_, err := g.userClient.Delete(context.Background(), &user_proto.UID{UID: ID})
	if err != nil {
		return errors.Wrap(err, "user client - Delete")
	}
	return nil
}

func (g UserClientGRPC) GetByID(ID uint64) (*models.User, error) {
	protoUser, err := g.userClient.GetByID(context.Background(), &user_proto.UID{UID: ID})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetByID")
	}
	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) GetByEmail(email string) (*models.User, error) {
	protoUser, err := g.userClient.GetByEmail(context.Background(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetByEmail")
	}
	return utils.UserModelByProto(protoUser), nil
}

func (g UserClientGRPC) GetInfo(ID uint64) (*models.UserInfo, error) {
	protoInfo, err := g.userClient.GetInfo(context.Background(), &user_proto.UID{UID: ID})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) GetInfoByEmail(email string) (*models.UserInfo, error) {
	fmt.Println()
	protoInfo, err := g.userClient.GetInfoByEmail(context.Background(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) EditInfo(ID uint64, info *models.UserInfo) (*models.UserInfo, error) {
	protoInfo, err := g.userClient.EditInfo(context.Background(), &user_proto.EditInfoParams{UID: ID, EditInfo: utils.ProtoByInfoModel(info)})
	if err != nil {
		return nil, errors.Wrap(err, "user client - EditInfo")
	}
	return utils.InfoModelByProto(protoInfo), nil
}

func (g UserClientGRPC) EditAvatar(ID uint64, newAvatar *models.Image, isCustom bool) error {
	_, err := g.userClient.EditAvatar(context.Background(), &user_proto.EditAvatarParams{
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
	protoImg, err := g.userClient.GetAvatar(context.Background(), &user_proto.Email{Email: email})
	if err != nil {
		return nil, errors.Wrap(err, "user client - GetAvatar")
	}
	return utils.ImageModelByProto(protoImg), nil
}

func (g UserClientGRPC) EditPw(ID uint64, form *models.EditPasswordRequest) error {
	_, err := g.userClient.EditPw(context.Background(), utils.ProtoByEditPasswordModel(form, ID))
	if err != nil {
		return errors.Wrap(err, "user client - EditPw")
	}
	return nil
}
