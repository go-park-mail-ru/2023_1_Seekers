package server

import (
	"context"
	"fmt"
	_user "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/user/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type UserServerGRPC struct {
	user_proto.UnimplementedUserServiceServer

	grpcServer *grpc.Server
	userUC     _user.UseCaseI
}

func NewUserServerGRPC(g *grpc.Server, uUC _user.UseCaseI) *UserServerGRPC {
	return &UserServerGRPC{
		grpcServer: g,
		userUC:     uUC,
	}
}

func (g *UserServerGRPC) Start(url string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return err
	}

	user_proto.RegisterUserServiceServer(g.grpcServer, g)
	return g.grpcServer.Serve(lis)
}

func (g *UserServerGRPC) Create(ctx context.Context, protoUser *user_proto.User) (*user_proto.User, error) {
	fmt.Println(utils.UserModelByProto(protoUser))
	user, err := g.userUC.Create(utils.UserModelByProto(protoUser))
	if err != nil {
		return nil, errors.Wrap(err, "user server - Create")
	}

	fmt.Println(utils.ProtoByUserModel(user))
	return utils.ProtoByUserModel(user), nil
}

func (g *UserServerGRPC) Delete(ctx context.Context, protoUID *user_proto.UID) (*user_proto.Nothing, error) {
	err := g.userUC.Delete(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "user server - Delete")
	}

	return &user_proto.Nothing{}, nil
}

func (g *UserServerGRPC) GetByID(ctx context.Context, protoUID *user_proto.UID) (*user_proto.User, error) {
	user, err := g.userUC.GetByID(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "user server - GetByID")
	}
	return utils.ProtoByUserModel(user), nil
}

func (g *UserServerGRPC) GetByEmail(ctx context.Context, protoEmail *user_proto.Email) (*user_proto.User, error) {
	user, err := g.userUC.GetByEmail(protoEmail.Email)
	if err != nil {
		return nil, errors.Wrap(err, "user server - GetByEmail")
	}
	return utils.ProtoByUserModel(user), nil
}

func (g *UserServerGRPC) GetInfo(ctx context.Context, protoUID *user_proto.UID) (*user_proto.UserInfo, error) {
	info, err := g.userUC.GetInfo(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "user server - GetInfo")
	}

	return utils.ProtoByInfoModel(info), nil
}

func (g *UserServerGRPC) GetInfoByEmail(ctx context.Context, protoEmail *user_proto.Email) (*user_proto.UserInfo, error) {
	info, err := g.userUC.GetInfoByEmail(protoEmail.Email)
	if err != nil {
		return nil, errors.Wrap(err, "user server - GetInfo")
	}

	return utils.ProtoByInfoModel(info), nil
}

func (g *UserServerGRPC) EditInfo(ctx context.Context, protoEditInfo *user_proto.EditInfoParams) (*user_proto.UserInfo, error) {
	info, err := g.userUC.EditInfo(protoEditInfo.UID, utils.InfoModelByProto(protoEditInfo.EditInfo))
	if err != nil {
		return nil, errors.Wrap(err, "user server - EditInfo")
	}
	return utils.ProtoByInfoModel(info), err
}

func (g *UserServerGRPC) EditAvatar(ctx context.Context, protoEditAvatar *user_proto.EditAvatarParams) (*user_proto.Nothing, error) {
	err := g.userUC.EditAvatar(protoEditAvatar.UID, utils.ImageModelByProto(protoEditAvatar.NewImage), protoEditAvatar.IsCustom)
	if err != nil {
		return nil, errors.Wrap(err, "user server - EditAvatar")
	}

	return &user_proto.Nothing{}, nil
}

func (g *UserServerGRPC) GetAvatar(ctx context.Context, protoEmail *user_proto.Email) (*user_proto.Image, error) {
	img, err := g.userUC.GetAvatar(protoEmail.Email)
	if err != nil {
		return nil, errors.Wrap(err, "user server - GetAvatar")
	}
	return utils.ProtoByImageModel(img), nil
}

func (g *UserServerGRPC) EditPw(ctx context.Context, protoEditPw *user_proto.EditPasswordParams) (*user_proto.Nothing, error) {
	err := g.userUC.EditPw(protoEditPw.UID, utils.EditPasswordModelByProto(protoEditPw))
	if err != nil {
		return nil, errors.Wrap(err, "user server - EditPw")
	}

	return &user_proto.Nothing{}, nil
}
