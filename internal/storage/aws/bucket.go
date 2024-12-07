package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Config struct {
	Region string `env:"AWS_REGION"`
	Bucket string `env:"AWS_BUCKET_NAME"`
}

type Repository struct {
	config   *Config
	s3Client s3iface.S3API
}

func NewRepository(config *Config) (*Repository, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(config.Region)}, nil)
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)

	return &Repository{config: config, s3Client: s3Client}, nil
}

func (r *Repository) DownloadFileFromBucket(ctx context.Context, fileName string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(r.config.Bucket),
		Key:    aws.String(fileName),
	}

	resp, err := r.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
