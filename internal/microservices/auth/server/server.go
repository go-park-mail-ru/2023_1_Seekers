package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/auth/utils"
	pkgGrpc "github.com/go-park-mail-ru/2023_1_Seekers/pkg/grpc"
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
		return nil, pkgGrpc.HandleError(ctx, err)
	}
	return utils.ProtoBySessionModel(session), nil
}

func (g *AuthServerGRPC) DeleteSession(ctx context.Context, protoSessionID *auth_proto.SessionId) (*auth_proto.Nothing, error) {
	err := g.authUC.DeleteSession(protoSessionID.Value)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return &auth_proto.Nothing{}, nil
}

func (g *AuthServerGRPC) GetSession(ctx context.Context, protoSessionID *auth_proto.SessionId) (*auth_proto.Session, error) {
	session, err := g.authUC.GetSession(protoSessionID.Value)
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}
	return utils.ProtoBySessionModel(session), nil
}

func (g *AuthServerGRPC) SignIn(ctx context.Context, protoFormLogin *auth_proto.FormLogin) (*auth_proto.AuthResponse, error) {
	info, session, err := g.authUC.SignIn(utils.LoginFormModelByProto(protoFormLogin))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoAuthResponseByInfoNSession(info, session), nil
}

func (g *AuthServerGRPC) SignUp(ctx context.Context, protoFormSignup *auth_proto.FormSignup) (*auth_proto.AuthResponse, error) {
	info, session, err := g.authUC.SignUp(utils.SignupFormModelByProto(protoFormSignup))
	if err != nil {
		return nil, pkgGrpc.HandleError(ctx, err)
	}

	return utils.ProtoAuthResponseByInfoNSession(info, session), nil
}
