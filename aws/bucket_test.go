package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	pngBytes  = append([]byte("\x89PNG\r\n\x1a\n"), bytes.Repeat([]byte{1}, 600)...)
	jpegBytes = append([]byte("\xff\xd8\xff\xe0"), bytes.Repeat([]byte{1}, 600)...)
	webpBytes = append([]byte("RIFF\x00\x00\x00\x00WEBPVP8 "), bytes.Repeat([]byte{1}, 600)...)
)

type fakeFile struct {
	*bytes.Reader
	closed bool
}

func (f *fakeFile) Close() error {
	f.closed = true
	return nil
}

func newFile(content []byte) (*fakeFile, *multipart.FileHeader) {
	return &fakeFile{Reader: bytes.NewReader(content)}, &multipart.FileHeader{Size: int64(len(content))}
}

type fakeS3 struct {
	putInputs    []*s3.PutObjectInput
	putBodies    [][]byte
	headErr      error
	listInputs   []*s3.ListObjectsV2Input
	listPages    []*s3.ListObjectsV2Output
	deleteInputs []*s3.DeleteObjectInput
}

func (f *fakeS3) PutObject(_ context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	body, _ := io.ReadAll(params.Body)
	f.putInputs = append(f.putInputs, params)
	f.putBodies = append(f.putBodies, body)
	return &s3.PutObjectOutput{}, nil
}

func (f *fakeS3) HeadObject(_ context.Context, _ *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if f.headErr != nil {
		return nil, f.headErr
	}
	return &s3.HeadObjectOutput{ContentLength: aws.Int64(42), LastModified: aws.Time(time.Unix(0, 0))}, nil
}

func (f *fakeS3) ListObjectsV2(_ context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	f.listInputs = append(f.listInputs, params)
	page := f.listPages[len(f.listInputs)-1]
	return page, nil
}

func (f *fakeS3) DeleteObject(_ context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	f.deleteInputs = append(f.deleteInputs, params)
	return &s3.DeleteObjectOutput{}, nil
}

// newTestBucket uses a real presign client: presigning is local and needs no network.
func newTestBucket(fake *fakeS3, now *time.Time) *awsBucket {
	client := s3.New(s3.Options{
		Region: "us-west-2",
		Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "AKIDTEST", SecretAccessKey: "secret"}, nil
		}),
	})
	return &awsBucket{
		s3Client:  fake,
		presigner: s3.NewPresignClient(client),
		signer: v4.NewSigner(func(so *v4.SignerOptions) {
			so.DisableURIPathEscaping = true
		}),
		bucketName: "olm-images",
		expiration: defaultURLExpiration,
		now:        func() time.Time { return *now },
	}
}

func windowStart() time.Time {
	return time.Date(2026, 11, 22, 15, 30, 0, 0, time.UTC).Truncate(12 * time.Hour)
}

func TestUploadSavesInsideFolderWithDetectedContentType(t *testing.T) {
	fake := &fakeS3{}
	now := windowStart().Add(time.Hour)
	bucket := newTestBucket(fake, &now)
	file, header := newFile(pngBytes)
	header.Header = map[string][]string{"Content-Type": {"text/html"}}

	object, err := bucket.Upload(context.Background(), &UploadParams{Folder: " /banners/home/ ", File: file, FileHeader: header})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(object.Key, "banners/home/") || !strings.HasSuffix(object.Key, ".png") {
		t.Errorf("unexpected key %q", object.Key)
	}
	input := fake.putInputs[0]
	if aws.ToString(input.Key) != object.Key || aws.ToString(input.Bucket) != "olm-images" {
		t.Errorf("unexpected put input key=%q bucket=%q", aws.ToString(input.Key), aws.ToString(input.Bucket))
	}
	if aws.ToString(input.ContentType) != "image/png" {
		t.Errorf("content type must come from the content, got %q", aws.ToString(input.ContentType))
	}
	if !bytes.Equal(fake.putBodies[0], pngBytes) {
		t.Error("the whole file must be uploaded from the start")
	}
	if !file.closed {
		t.Error("the file must be closed")
	}

	signed, _ := url.Parse(object.URL)
	query := signed.Query()
	if signed.Scheme != "https" || signed.Host != "olm-images.s3.us-west-2.amazonaws.com" || signed.Path != "/"+object.Key {
		t.Errorf("unexpected URL %s", object.URL)
	}
	if query.Get("X-Amz-Expires") != "86400" || query.Get("X-Amz-Signature") == "" {
		t.Errorf("URL is not presigned for 24h: %s", object.URL)
	}
	if !object.ExpiresAt.Equal(windowStart().Add(24 * time.Hour)) {
		t.Errorf("unexpected expiration %s", object.ExpiresAt)
	}
}

func TestUploadAcceptsJpegAndWebp(t *testing.T) {
	cases := map[string]struct {
		content     []byte
		contentType string
		extension   string
	}{
		"jpeg": {jpegBytes, "image/jpeg", ".jpg"},
		"webp": {webpBytes, "image/webp", ".webp"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fake := &fakeS3{}
			now := windowStart()
			file, header := newFile(tc.content)

			object, err := newTestBucket(fake, &now).Upload(context.Background(), &UploadParams{Folder: "avatar", File: file, FileHeader: header})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.HasSuffix(object.Key, tc.extension) || aws.ToString(fake.putInputs[0].ContentType) != tc.contentType {
				t.Errorf("unexpected key %q or content type %q", object.Key, aws.ToString(fake.putInputs[0].ContentType))
			}
		})
	}
}

func TestUploadRejectsInvalidInput(t *testing.T) {
	cases := map[string]struct {
		folder  string
		content []byte
		size    int64
		want    error
	}{
		"html":        {"banners", []byte("<html><body>hi</body></html>"), 0, ErrInvalidContentType},
		"svg":         {"banners", []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"></svg>`), 0, ErrInvalidContentType},
		"too large":   {"banners", pngBytes, multipartMaxLength + 1, ErrFileTooLarge},
		"no folder":   {"", pngBytes, 0, ErrInvalidFolder},
		"only slash":  {"/", pngBytes, 0, ErrInvalidFolder},
		"parent":      {"../banners", pngBytes, 0, ErrInvalidFolder},
		"empty level": {"banners//home", pngBytes, 0, ErrInvalidFolder},
		"dot level":   {"banners/./home", pngBytes, 0, ErrInvalidFolder},
		"space":       {"my banners", pngBytes, 0, ErrInvalidFolder},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fake := &fakeS3{}
			now := windowStart()
			file, header := newFile(tc.content)
			if tc.size > 0 {
				header.Size = tc.size
			}

			_, err := newTestBucket(fake, &now).Upload(context.Background(), &UploadParams{Folder: tc.folder, File: file, FileHeader: header})
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
			if len(fake.putInputs) != 0 {
				t.Error("nothing must be uploaded")
			}
			if !file.closed {
				t.Error("the file must be closed")
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	now := windowStart()

	object, err := newTestBucket(&fakeS3{}, &now).GetFile(context.Background(), "banners/a.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if object.Key != "banners/a.png" || object.Size != 42 || !strings.Contains(object.URL, "X-Amz-Signature=") {
		t.Errorf("unexpected object %+v", object)
	}

	notFound := &fakeS3{headErr: fmt.Errorf("operation error S3: HeadObject: %w", &types.NotFound{})}
	if _, err = newTestBucket(notFound, &now).GetFile(context.Background(), "banners/missing.png"); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}

	for _, key := range []string{"", "/banners/a.png", "banners/", "../a.png"} {
		if _, err = newTestBucket(&fakeS3{}, &now).GetFile(context.Background(), key); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("key %q: want ErrInvalidKey, got %v", key, err)
		}
	}
}

func TestListFolderReadsEveryPage(t *testing.T) {
	fake := &fakeS3{listPages: []*s3.ListObjectsV2Output{
		{
			Contents: []types.Object{
				{Key: aws.String("banners/"), Size: aws.Int64(0)},
				{Key: aws.String("banners/a.png"), Size: aws.Int64(10)},
			},
			IsTruncated:           aws.Bool(true),
			NextContinuationToken: aws.String("next"),
		},
		{
			Contents:    []types.Object{{Key: aws.String("banners/b.jpg"), Size: aws.Int64(20)}},
			IsTruncated: aws.Bool(false),
		},
	}}
	now := windowStart()

	objects, err := newTestBucket(fake, &now).ListFolder(context.Background(), "banners/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(objects) != 2 || objects[0].Key != "banners/a.png" || objects[1].Key != "banners/b.jpg" {
		t.Fatalf("unexpected objects %+v", objects)
	}
	for _, object := range objects {
		if !strings.Contains(object.URL, "X-Amz-Signature=") {
			t.Errorf("object %s is not presigned", object.Key)
		}
	}
	first := fake.listInputs[0]
	if aws.ToString(first.Prefix) != "banners/" || aws.ToString(first.Delimiter) != "/" {
		t.Errorf("unexpected prefix %q or delimiter %q", aws.ToString(first.Prefix), aws.ToString(first.Delimiter))
	}
	if aws.ToString(fake.listInputs[1].ContinuationToken) != "next" {
		t.Error("the second page must use the continuation token")
	}
}

func TestListFolderEmptyReturnsEmptySlice(t *testing.T) {
	fake := &fakeS3{listPages: []*s3.ListObjectsV2Output{{IsTruncated: aws.Bool(false)}}}
	now := windowStart()

	objects, err := newTestBucket(fake, &now).ListFolder(context.Background(), "avatar")
	if err != nil || objects == nil || len(objects) != 0 {
		t.Errorf("want an empty slice, got %v, %v", objects, err)
	}

	if _, err = newTestBucket(&fakeS3{}, &now).ListFolder(context.Background(), ".."); !errors.Is(err, ErrInvalidFolder) {
		t.Errorf("want ErrInvalidFolder, got %v", err)
	}
}

func TestURLIsStableInsideTheSigningWindow(t *testing.T) {
	now := windowStart().Add(time.Hour)
	bucket := newTestBucket(&fakeS3{}, &now)
	get := func() *Object {
		object, err := bucket.GetFile(context.Background(), "avatar/a.png")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return object
	}

	first := get()
	now = windowStart().Add(12*time.Hour - time.Second)
	sameWindow := get()
	now = windowStart().Add(12 * time.Hour)
	nextWindow := get()

	if first.URL != sameWindow.URL {
		t.Errorf("URL changed inside the window:\n%s\n%s", first.URL, sameWindow.URL)
	}
	if nextWindow.URL == first.URL {
		t.Error("URL must change in the next window")
	}
	if !nextWindow.ExpiresAt.Equal(windowStart().Add(36 * time.Hour)) {
		t.Errorf("unexpected expiration %s", nextWindow.ExpiresAt)
	}
	signed, _ := url.Parse(first.URL)
	if signed.Query().Get("X-Amz-Date") != windowStart().Format("20060102T150405Z") {
		t.Errorf("URL must be signed at the window start, got %s", signed.Query().Get("X-Amz-Date"))
	}
}

func TestDeleteFile(t *testing.T) {
	fake := &fakeS3{}
	now := windowStart()
	bucket := newTestBucket(fake, &now)

	if err := bucket.DeleteFile(context.Background(), "avatar/a.png"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if aws.ToString(fake.deleteInputs[0].Key) != "avatar/a.png" {
		t.Errorf("unexpected key %q", aws.ToString(fake.deleteInputs[0].Key))
	}

	if err := bucket.DeleteFile(context.Background(), "avatar/"); !errors.Is(err, ErrInvalidKey) {
		t.Errorf("want ErrInvalidKey, got %v", err)
	}
	if len(fake.deleteInputs) != 1 {
		t.Error("an invalid key must not reach the bucket")
	}
}

func TestNewAwsBucketValidatesConfig(t *testing.T) {
	valid := BucketConfig{AccessKey: "a", SecretKey: "s", BucketName: "olm-images", Region: "us-west-2"}

	service, err := NewAwsBucket(&valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service.(*awsBucket).expiration != defaultURLExpiration {
		t.Error("expiration must default to 24h")
	}

	tooLong := valid
	tooLong.URLExpiration = maxURLExpiration + time.Second
	if _, err = NewAwsBucket(&tooLong); err == nil {
		t.Error("expiration above 7 days must be refused")
	}

	noBucket := valid
	noBucket.BucketName = ""
	if _, err = NewAwsBucket(&noBucket); err == nil {
		t.Error("an empty bucket name must be refused")
	}
}
