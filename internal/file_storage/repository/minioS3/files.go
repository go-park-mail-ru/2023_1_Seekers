package minioS3

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type fileDB struct {
	uploaderS3   *s3manager.Uploader
	downloaderS3 *s3manager.Downloader
}

func New() file_storage.RepoI {
	s, err := session.NewSession(
		&aws.Config{
			Endpoint:         aws.String(config.S3Endpoint),
			Region:           aws.String(config.S3Region),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv(config.S3AccessKeyEnv),
				os.Getenv(config.S3ASecretKeyEnv),
				"",
			),
		},
	)
	if err != nil {
		logrus.Fatalf("Failed create S3 session : %v", err)
	}

	db := &fileDB{
		s3manager.NewUploader(s),
		s3manager.NewDownloader(s),
	}

	dat, err := os.ReadFile(config.DefaultAvatarDir + config.DefaultAvatar)
	if err != nil {
		logrus.Fatalf("Failed init default avatar: %v", err)
	}

	err = db.Upload(&models.S3File{
		Bucket: config.S3AvatarBucket,
		Name:   config.DefaultAvatar,
		Data:   dat,
	})

	if err != nil {
		logrus.Fatalf("Failed init default avatar: %v", err)
	}

	return db
}

func (fDB *fileDB) Get(bName, fName string) (*models.S3File, error) {
	objInput := &awsS3.GetObjectInput{
		Bucket: &bName,
		Key:    &fName,
	}
	buf := &aws.WriteAtBuffer{}
	numBytes, err := fDB.downloaderS3.Download(buf, objInput)
	if err != nil {
		return nil, err
	}
	if numBytes < 1 {
		return nil, errors.New("written empty file")
	}

	return &models.S3File{
		Bucket: bName,
		Name:   fName,
		Data:   buf.Bytes(),
	}, nil
}

func (fDB *fileDB) Upload(file *models.S3File) error {
	uInput := &s3manager.UploadInput{
		Bucket:      aws.String(file.Bucket),
		Key:         aws.String(file.Name),
		Body:        bytes.NewReader(file.Data),
		ContentType: aws.String(http.DetectContentType(file.Data)),
		//ACL:         aws.String("public-read"),
	}

	_, err := fDB.uploaderS3.Upload(uInput)
	if err != nil {
		return err
	}

	return nil
}
