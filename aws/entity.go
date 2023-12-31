package aws

import (
	"context"
	"mime/multipart"
)

type UploadFileParams struct {
	ProductID  string
	File       multipart.File
	FileHeader *multipart.FileHeader
	Ctx        context.Context
}

type BucketConfig struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	Region     string
}
