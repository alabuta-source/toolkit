package aws

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"
)

var (
	multipartMaxLength   int64 = 4 << 20 // 4MB
	acceptedContentTypes       = []string{"image/png", "image/jpeg"}
)

func buildPublicURL(fileID string, fileExtension string, bucket string) string {
	return fmt.Sprintf(
		"http://%s.s3-website-us-west-2.amazonaws.com/%s.%s",
		bucket,
		fileID,
		fileExtension,
	)
}

func hasValidContentType(fileHeader *multipart.FileHeader) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	for _, ct := range acceptedContentTypes {
		if ct == contentType {
			return true
		}
	}
	return false
}

func getFileType(fileHeader *multipart.FileHeader) string {
	contentType := fileHeader.Header.Get("Content-Type")
	return strings.ReplaceAll(contentType, "image/", "")
}

func generateUUID() string {
	return uuid.New().String()
}
