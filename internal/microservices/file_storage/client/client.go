package client

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/utils"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type FStorageClientGRPC struct {
	fStorageClient fstorage_proto.FileStorageServiceClient
}

func NewFstorageClientGRPC(cc *grpc.ClientConn) file_storage.UseCaseI {
	return &FStorageClientGRPC{
		fStorageClient: fstorage_proto.NewFileStorageServiceClient(cc),
	}
}

func (f *FStorageClientGRPC) Get(bName, fName string) (*models.S3File, error) {
	file, err := f.fStorageClient.Get(context.TODO(), utils.ProtoGetParams(bName, fName))
	if err != nil {
		return nil, errors.Wrap(err, "fStorage client - Get")
	}
	return utils.FileModelByProto(file), nil
}

func (f *FStorageClientGRPC) Upload(file *models.S3File) error {
	_, err := f.fStorageClient.Upload(context.TODO(), utils.ProtoFileByModel(file))
	if err != nil {
		return errors.Wrap(err, "fStorage client - Upload")
	}
	return nil
}
