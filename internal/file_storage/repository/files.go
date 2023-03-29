package repository

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-park-mail-ru/2023_1_Seekers/cmd/config"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/file_storage"
	"github.com/go-park-mail-ru/2023_1_Seekers/internal/models"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	pkgErrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type fileDB struct {
	uploaderS3   *s3manager.Uploader
	downloaderS3 *s3manager.Downloader
}

func New(s *session.Session) file_storage.RepoI {
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
		var awsErr awserr.Error
		if pkgErrors.As(err, &awsErr) {
			switch awsErr.Code() {
			case awsS3.ErrCodeNoSuchBucket:
				return nil, pkgErrors.WithMessagef(errors.ErrNoBucket,
					"bucket [%s], key [%s], error [%s]", bName, fName, err.Error())
			case awsS3.ErrCodeNoSuchKey:
				return nil, pkgErrors.WithMessagef(errors.ErrNoKey,
					"bucket [%s], key [%s], error [%s]", bName, fName, err.Error())
			}
		}
		return nil, pkgErrors.WithMessagef(errors.ErrGetFile, err.Error())
	}
	if numBytes < 1 {
		return nil, pkgErrors.WithMessagef(errors.ErrGetFile,
			"empty file : bucket [%s], key [%s], error [%s]", bName, fName, err.Error())
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
		return pkgErrors.WithMessagef(errors.ErrInternal,
			"upload file: bucket [%s], key [%s], error [%s]", file.Bucket, file.Name, err)
	}

	return nil
}
