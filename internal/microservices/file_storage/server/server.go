package server

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type FStorageServerGRPC struct {
	fstorage_proto.UnimplementedFileStorageServiceServer

	grpcServer *grpc.Server
	fStorageUC file_storage.UseCaseI
}

func NewFStorageServerGRPC(grpcSrv *grpc.Server, fSUC file_storage.UseCaseI) *FStorageServerGRPC {
	return &FStorageServerGRPC{

		grpcServer: grpcSrv,
		fStorageUC: fSUC,
	}
}

func (g *FStorageServerGRPC) Start(url string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return err
	}
	fstorage_proto.RegisterFileStorageServiceServer(g.grpcServer, g)
	return g.grpcServer.Serve(lis)
}

func (g *FStorageServerGRPC) Get(ctx context.Context, params *fstorage_proto.GetFileParams) (*fstorage_proto.File, error) {
	file, err := g.fStorageUC.Get(params.BucketName, params.FileName)
	if err != nil {
		return nil, errors.Wrap(err, "file storage - get")
	}
	return utils.ProtoFileByModel(file), nil
}

func (g *FStorageServerGRPC) Upload(ctx context.Context, protoFile *fstorage_proto.File) (*fstorage_proto.Nothing, error) {
	if err := g.fStorageUC.Upload(utils.FileModelByProto(protoFile)); err != nil {
		return nil, errors.Wrap(err, "file storage - upload")
	}
	return nil, nil
}
