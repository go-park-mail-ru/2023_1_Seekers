package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	_fStorageRepo "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/repository/S3"
	_fStorageServer "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/server"
	_fStorageUCase "github.com/go-park-mail-ru/2023_1_Seekers/internal/microservices/file_storage/usecase"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/connectors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
)

func main() {
	endpoint := aws.String(config.S3Endpoint)
	region := aws.String(config.S3Region)
	disableSSL := aws.Bool(true)
	s3ForcePathStyle := aws.Bool(true)
	creds := credentials.NewStaticCredentials(
		os.Getenv(config.S3AccessKeyEnv),
		os.Getenv(config.S3ASecretKeyEnv),
		"",
	)
	s3Session, err := connectors.NewS3(endpoint, region, disableSSL, s3ForcePathStyle, creds)
	if err != nil {
		log.Fatalf("Failecd create S3 session : %v", err)
	}

	fStorageRepo := _fStorageRepo.New(s3Session)
	fStorageUC := _fStorageUCase.New(fStorageRepo)

	grpcServer := grpc.NewServer()
	fStorageGRPCServer := _fStorageServer.NewFStorageServerGRPC(grpcServer, fStorageUC)
	log.Info("file storage server started")
	err = fStorageGRPCServer.Start(":" + config.FileServiceGRPCPort)
	if err != nil {
		log.Fatal(err)
	}
}
