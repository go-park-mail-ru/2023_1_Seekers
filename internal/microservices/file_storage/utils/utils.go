package utils

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/proto"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
)

func FileModelByProto(protoFile *fstorage_proto.File) *models.S3File {
	return &models.S3File{
		Bucket: protoFile.BucketName,
		Name:   protoFile.FileName,
		Data:   protoFile.Data,
	}
}

func ProtoFileByModel(file *models.S3File) *fstorage_proto.File {
	return &fstorage_proto.File{
		BucketName: file.Bucket,
		FileName:   file.Name,
		Data:       file.Data,
	}
}

func ProtoGetParams(bName, fName string) *fstorage_proto.GetFileParams {
	return &fstorage_proto.GetFileParams{
		BucketName: bName,
		FileName:   fName,
	}
}
