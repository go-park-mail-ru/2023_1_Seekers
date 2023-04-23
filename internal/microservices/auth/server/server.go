package server

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type AuthServerGRPC struct {
	auth_proto.UnimplementedAuthServiceServer

	grpcServer *grpc.Server
	authUC     auth.UseCaseI
}

func NewAuthServerGRPC(g *grpc.Server, aUC auth.UseCaseI) *AuthServerGRPC {
	return &AuthServerGRPC{
		grpcServer: g,
		authUC:     aUC,
	}
}

func (g *AuthServerGRPC) Start(url string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return err
	}
	auth_proto.RegisterAuthServiceServer(g.grpcServer, g)
	return g.grpcServer.Serve(lis)
}

func (g *AuthServerGRPC) CreateSession(ctx context.Context, protoUID *auth_proto.UID) (*auth_proto.Session, error) {
	session, err := g.authUC.CreateSession(protoUID.UID)
	if err != nil {
		return nil, errors.Wrap(err, "auth server - CreateSession")
	}
	return utils.ProtoBySessionModel(session), nil
}

func (g *AuthServerGRPC) DeleteSession(_ context.Context, protoSessionID *auth_proto.SessionId) (*auth_proto.Nothing, error) {
	err := g.authUC.DeleteSession(protoSessionID.Value)
	if err != nil {
		return nil, errors.Wrap(err, "auth server - DeleteSession")
	}

	return &auth_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) GetSession(_ context.Context, protoSessionID *auth_proto.SessionId) (*auth_proto.Session, error) {
	session, err := g.authUC.GetSession(protoSessionID.Value)
	if err != nil {
		return nil, errors.Wrap(err, "auth server - GetSession")
	}
	return utils.ProtoBySessionModel(session), nil
}

func (g *AuthServerGRPC) SignIn(_ context.Context, protoFormLogin *auth_proto.FormLogin) (*auth_proto.AuthResponse, error) {
	info, session, err := g.authUC.SignIn(utils.LoginFormModelByProto(protoFormLogin))
	if err != nil {
		return nil, errors.Wrap(err, "auth server - SignIn")
	}

	return utils.ProtoAuthResponseByInfoNSession(info, session), nil
}

func (g *AuthServerGRPC) SignUp(_ context.Context, protoFormSignup *auth_proto.FormSignup) (*auth_proto.AuthResponse, error) {
	fmt.Println(utils.SignupFormModelByProto(protoFormSignup))
	info, session, err := g.authUC.SignUp(utils.SignupFormModelByProto(protoFormSignup))
	if err != nil {
		return nil, errors.Wrap(err, "auth server - SignUp")
	}

	return utils.ProtoAuthResponseByInfoNSession(info, session), nil
}
