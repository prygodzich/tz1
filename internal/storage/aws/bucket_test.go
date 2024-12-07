package aws

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockS3Client is a mock of S3 client
type MockS3Client struct {
	s3iface.S3API
	getObjectOutput *s3.GetObjectOutput
	getObjectError  error
}

func (m *MockS3Client) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	return m.getObjectOutput, m.getObjectError
}

func TestRepository_DownloadFileFromBucket(t *testing.T) {
	tests := []struct {
		name          string
		fileName      string
		mockS3Client  *MockS3Client
		expectedData  []byte
		expectedError error
	}{
		{
			name:     "Successful download",
			fileName: "test.txt",
			mockS3Client: &MockS3Client{
				getObjectOutput: &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test data"))),
				},
				getObjectError: nil,
			},
			expectedData:  []byte("test data"),
			expectedError: nil,
		},
		{
			name:     "File not found",
			fileName: "nonexistent.txt",
			mockS3Client: &MockS3Client{
				getObjectOutput: nil,
				getObjectError:  awserr.New(s3.ErrCodeNoSuchKey, "The specified key does not exist.", nil),
			},
			expectedData:  nil,
			expectedError: awserr.New(s3.ErrCodeNoSuchKey, "The specified key does not exist.", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Region: "us-west-2",
				Bucket: "test-bucket",
			}
			repo := &Repository{
				config:   config,
				s3Client: tt.mockS3Client, // Используем наш мок-клиент
			}

			ctx := context.Background()
			data, err := repo.DownloadFileFromBucket(ctx, tt.fileName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				awsErr, ok := err.(awserr.Error)
				require.True(t, ok, "Expected AWS error")
				assert.Equal(t, tt.expectedError.(awserr.Error).Code(), awsErr.Code())
				assert.Equal(t, tt.expectedError.(awserr.Error).Message(), awsErr.Message())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}
