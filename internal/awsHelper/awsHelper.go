package awsHelper

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Helper struct {
	S3     *s3.Client
	bucket string
}

func New(accessKeyId, accessKeySecret, accountId string, bucket string) *S3Helper {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return &S3Helper{
		S3:     client,
		bucket: bucket,
	}
}

func (h *S3Helper) GetObject(key string) (*s3.GetObjectOutput, error) {
	result, err := h.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func IsNotFoundError(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "NoSuchKey"
	}
	return false
}
