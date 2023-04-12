package connectors

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewS3(endpoint, region *string, disableSSL, s3ForcePathStyle *bool, credentials *credentials.Credentials) (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Endpoint:         endpoint,
			Region:           region,
			DisableSSL:       disableSSL,
			S3ForcePathStyle: s3ForcePathStyle,
			Credentials:      credentials,
		},
	)
}
