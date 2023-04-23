package client

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type AuthClientGRPC struct {
	authClient auth_proto.AuthServiceClient
}

func NewAuthClientGRPC(cc *grpc.ClientConn) auth.UseCaseI {
	return &AuthClientGRPC{
		authClient: auth_proto.NewAuthServiceClient(cc),
	}
}

func (g AuthClientGRPC) SignIn(form *models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	authResp, err := g.authClient.SignIn(context.Background(), utils.ProtoByLoginFormModel(form))
	if err != nil {
		return nil, nil, errors.Wrap(err, "auth client - SignIn")
	}

	return utils.AuthResponseModelByProto(authResp), utils.SessionModelByProtoAuthResponse(authResp), nil
}

func (g AuthClientGRPC) SignUp(form *models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	fmt.Println("---------")
	fmt.Println(utils.ProtoBySignupFormModel(form))
	authResp, err := g.authClient.SignUp(context.Background(), utils.ProtoBySignupFormModel(form))
	if err != nil {
		return nil, nil, errors.Wrap(err, "auth client - SignUp")
	}

	return utils.AuthResponseModelByProto(authResp), utils.SessionModelByProtoAuthResponse(authResp), nil
}

func (g AuthClientGRPC) CreateSession(uID uint64) (*models.Session, error) {
	protoSession, err := g.authClient.CreateSession(context.Background(), &auth_proto.UID{UID: uID})
	if err != nil {
		return nil, errors.Wrap(err, "auth client - CreateSession")
	}

	return utils.SessionModelByProto(protoSession), nil
}

func (g AuthClientGRPC) DeleteSession(sessionID string) error {
	_, err := g.authClient.DeleteSession(context.Background(), &auth_proto.SessionId{Value: sessionID})
	if err != nil {
		return errors.Wrap(err, "auth client - DeleteSession")
	}

	return nil
}

func (g AuthClientGRPC) GetSession(sessionID string) (*models.Session, error) {
	protoSession, err := g.authClient.GetSession(context.Background(), &auth_proto.SessionId{Value: sessionID})
	if err != nil {
		return nil, errors.Wrap(err, "auth client - GetSession")
	}

	return utils.SessionModelByProto(protoSession), nil
}
