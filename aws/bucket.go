package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"mime/multipart"
)

type BucketService interface {
	UploadFile(params *UploadFileParams) (string, error)
	ListObjetFiles(ctx context.Context, productID string) ([]string, error)
}

type awsBucket struct {
	s3Client   *s3.Client
	bucketName string
}

func NewAwsBucket(configs *BucketConfig) BucketService {
	client := s3.New(s3.Options{
		Region: configs.Region,
		Credentials: aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				cred := aws.Credentials{
					AccessKeyID:     configs.AccessKey,
					SecretAccessKey: configs.SecretKey,
				}

				if !cred.HasKeys() {
					return aws.Credentials{}, errors.New("the keys are missing")
				}
				return cred, nil
			},
		),
	})
	return &awsBucket{
		s3Client:   client,
		bucketName: configs.BucketName,
	}
}

func (bucket *awsBucket) UploadFile(params *UploadFileParams) (string, error) {
	if !hasValidContentType(params.FileHeader) {
		return "", errors.New("invalid Content-Type, here is the valid list ['image/png', 'image/jpeg']")
	}

	if params.FileHeader.Size > multipartMaxLength {
		return "", fmt.Errorf("image too large, max len: %d [4MB]", multipartMaxLength)
	}

	defer func(File multipart.File) {
		err := File.Close()
		if err != nil {
			log.Println("[toolkit.UploadFile] error closing the multipart file")
		}
	}(params.File)

	fileType := getFileType(params.FileHeader)
	fileName := fmt.Sprintf("%s/%s.%s", params.ProductID, generateUUID(), fileType)

	_, err := bucket.s3Client.PutObject(params.Ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket.bucketName),
		Key:         aws.String(fileName),
		Body:        params.File,
		ContentType: aws.String(params.FileHeader.Header.Get(contentTypeKey)),
	})
	if err != nil {
		return "", err
	}
	return buildPublicURL(fileName, bucket.bucketName), nil
}

func (bucket *awsBucket) ListObjetFiles(ctx context.Context, productID string) ([]string, error) {
	result, err := bucket.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket.bucketName),
		Prefix: aws.String(fmt.Sprintf("%s/", productID)),
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list objects in bucket %s. Here's why: %v\n",
			bucket.bucketName,
			err,
		)
	}
	return buildAnyPublicURL(result.Contents, bucket.bucketName), err
}
