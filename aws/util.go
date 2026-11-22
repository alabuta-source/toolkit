package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/google/uuid"
)

const (
	multipartMaxLength   int64 = 4 << 20 // 4MB
	defaultURLExpiration       = 24 * time.Hour
	maxURLExpiration           = 7 * 24 * time.Hour // SigV4 limit
	sniffLength                = 512
)

var (
	ErrNotFound           = errors.New("file not found")
	ErrInvalidFolder      = errors.New("invalid folder, use letters, numbers, '-' or '_' separated by '/'")
	ErrInvalidKey         = errors.New("invalid file key")
	ErrInvalidContentType = errors.New("invalid Content-Type, here is the valid list ['image/png', 'image/jpeg', 'image/webp']")
	ErrFileTooLarge       = fmt.Errorf("image too large, max len: %d [4MB]", multipartMaxLength)
)

// acceptedContentTypes maps the type detected from the file content to the key extension.
var acceptedContentTypes = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/webp": "webp",
}

var folderPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+(/[a-zA-Z0-9_-]+)*$`)

func normalizeFolder(folder string) (string, error) {
	folder = strings.Trim(strings.TrimSpace(folder), "/")
	if !folderPattern.MatchString(folder) {
		return "", ErrInvalidFolder
	}
	return folder, nil
}

func validateKey(key string) error {
	if key == "" || strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") || strings.Contains(key, "..") {
		return ErrInvalidKey
	}
	return nil
}

// detectContentType reads the file content, ignoring the Content-Type sent by the client,
// and rewinds the file so it can be uploaded from the start.
func detectContentType(file multipart.File) (string, error) {
	buf := make([]byte, sniffLength)
	n, err := io.ReadFull(file, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("error reading the file: %w", err)
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("error rewinding the file: %w", err)
	}

	contentType := http.DetectContentType(buf[:n])
	if _, ok := acceptedContentTypes[contentType]; !ok {
		return "", ErrInvalidContentType
	}
	return contentType, nil
}

func buildKey(folder, contentType string) string {
	return fmt.Sprintf("%s/%s.%s", folder, uuid.New().String(), acceptedContentTypes[contentType])
}

// signingWindow is the period during which an object keeps the same presigned URL, so the app can
// cache the image by URL. Half of the expiration leaves at least half of it on every URL handed out.
func signingWindow(expiration time.Duration) time.Duration {
	return expiration / 2
}

// fixedTimePresigner signs with a given time instead of the current one, which keeps the URL stable
// inside a signing window.
type fixedTimePresigner struct {
	signer   *v4.Signer
	signedAt time.Time
}

func (p fixedTimePresigner) PresignHTTP(
	ctx context.Context, credentials aws.Credentials, r *http.Request,
	payloadHash string, service string, region string, _ time.Time,
	optFns ...func(*v4.SignerOptions),
) (string, http.Header, error) {
	return p.signer.PresignHTTP(ctx, credentials, r, payloadHash, service, region, p.signedAt, optFns...)
}

func brasiliaTime() time.Time {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.FixedZone("America/Sao_Paulo", -3*60*60)
	}
	return time.Now().In(loc)
}
