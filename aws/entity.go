package aws

import (
	"mime/multipart"
	"time"
)

type UploadParams struct {
	// Folder groups the image inside the bucket, example: "banners", "avatar" or "banners/home".
	Folder     string
	File       multipart.File
	FileHeader *multipart.FileHeader
}

// Object is an image stored in the bucket.
// Persist the Key: the URL is presigned, expires at ExpiresAt and must be generated again from the Key.
type Object struct {
	Key          string    `json:"key"`
	URL          string    `json:"url"`
	ExpiresAt    time.Time `json:"expires_at"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
}

type BucketConfig struct {
	AccessKey  string
	SecretKey  string
	BucketName string
	Region     string
	// URLExpiration is how long a presigned URL stays valid. Defaults to 24h, max 7 days.
	URLExpiration time.Duration
}
