package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"
)

var (
	multipartMaxLength   int64 = 4 << 20 // 4MB
	acceptedContentTypes       = []string{"image/png", "image/jpeg"}
	contentTypeKey             = "Content-Type"
)

func buildPublicURL(fileID string, bucket string) string {
	return fmt.Sprintf(
		"http://%s.s3-website-us-west-2.amazonaws.com/%s",
		bucket,
		fileID,
	)
}

func buildAnyPublicURL(contents []types.Object, bucketName string) []string {
	if len(contents) < 1 {
		return nil
	}
	urls := make([]string, 0, len(contents))
	for _, content := range contents {
		urls = append(urls, buildPublicURL(*content.Key, bucketName))
	}
	return urls
}

func hasValidContentType(fileHeader *multipart.FileHeader) bool {
	contentType := fileHeader.Header.Get(contentTypeKey)
	for _, ct := range acceptedContentTypes {
		if ct == contentType {
			return true
		}
	}
	return false
}

func getFileType(fileHeader *multipart.FileHeader) string {
	contentType := fileHeader.Header.Get(contentTypeKey)
	return strings.ReplaceAll(contentType, "image/", "")
}

func generateUUID() string {
	return uuid.New().String()
}
