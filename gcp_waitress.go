package toolkit

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

var (
	MultipartMaxLength   int64 = 4 << 20 // 4MB
	acceptedContentTypes       = map[string]bool{"image/png": true, "image/jpeg": true}
)

type GCPWaitressManager interface {
	// UploadFile saves a file to the bucket and returns the name of the file or an error
	// The object is the name of the file to save in the bucket
	// The prefix is used to create a folder in the bucket, if the prefix is empty, the file will be saved in the root of the bucket.
	UploadFile(file multipart.File, fileHeader *multipart.FileHeader, prefix string) (string, error)

	// ListFiles returns a list of files in the bucket
	// The prefix is used to filter the files, if the prefix is empty, all files will be returned.
	// Make sure to add the prefix to filter the files, example: "users"
	ListFiles(prefix string) ([]string, error)

	// DeleteFile use this method to delete a single file.
	DeleteFile(filename string) error
}

type gcpWaitress struct {
	client     *storage.Client
	ctx        context.Context
	bucket     *storage.BucketHandle
	bucketName string
}

// NewGCPWaitress creates a new GCP Waitress pointer to access the GCP Bucket
// The bucket is the name of the bucket to access in Google Cloud Storage
// The request is the http request to get the context, this is used to create the context for the GCP Client
// The gcpKey is the json key to access the bucket, this key can be created in the Google Cloud Console and should have the Storage Admin role
func NewGCPWaitress(bucketName string, request *http.Request, gcpKey *GCPBucketAuthJson) (GCPWaitressManager, error) {
	ctx := appengine.NewContext(request)
	bytes, er := json.Marshal(gcpKey)
	if er != nil {
		return nil, fmt.Errorf("error trying marshal de GCP key")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(bytes))
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	if !bucketExiste(ctx, bucket) {
		return nil, errors.New("bucket not founded")
	}

	return &gcpWaitress{
		client:     client,
		ctx:        ctx,
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

func (w *gcpWaitress) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, prefix string) (string, error) {
	if !w.hasValidContentType(fileHeader) {
		return "", errors.New("invalid Content-Type, here is the valid list ['image/png', 'image/jpeg']")
	}

	if fileHeader.Size > MultipartMaxLength {
		return "", fmt.Errorf("image too large, max len: %d [4MB]", MultipartMaxLength)
	}

	id := GenerateUUID()
	name := fmt.Sprintf("%s/%s", prefix, id)
	if prefix == "" {
		log.Printf("You're saving file:[%s] without prefix", name)
		name = id
	}

	wc := w.bucket.Object(name).NewWriter(w.ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return w.buildURL(wc.Name), nil
}

func (w *gcpWaitress) ListFiles(prefix string) ([]string, error) {
	resp := make([]string, 0)
	objects := w.bucket.Objects(w.ctx, &storage.Query{
		StartOffset: prefix + "/",
	})

	for {
		attrs, err := objects.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error loading file: %w", err)
		}
		resp = append(resp, attrs.Name)
	}
	return w.buildURLs(resp), nil
}

func (w *gcpWaitress) DeleteFile(fileUrl string) error {
	filename, er := w.objectNameFromUrl(fileUrl)
	if er != nil {
		return er
	}

	oHandle := w.bucket.Object(filename)
	attrs, err := oHandle.Attrs(w.ctx)
	if err != nil {
		return fmt.Errorf("error getting meta information about the object to delete it: %w", err)
	}

	if er = oHandle.
		If(storage.Conditions{GenerationMatch: attrs.Generation}).
		Delete(w.ctx); er != nil {
		return fmt.Errorf("error trying delete obj: %w", er)
	}
	return nil
}

func bucketExiste(ctx context.Context, b *storage.BucketHandle) bool {
	_, err := b.Attrs(ctx)
	return err == nil
}

func (w *gcpWaitress) buildURL(name string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", w.bucketName, name)
}

func (w *gcpWaitress) buildURLs(names []string) []string {
	resp := make([]string, 0, len(names))
	for _, name := range names {
		resp = append(resp, w.buildURL(name))
	}
	return resp
}

func (w *gcpWaitress) objectNameFromUrl(imgUrl string) (string, error) {
	if imgUrl == "" {
		return "", errors.New("the object url should not be empty")
	}

	urlPath, err := url.Parse(imgUrl)
	if err != nil {
		return "", fmt.Errorf("unable to parse the url: %s", err.Error())
	}
	return removeBucketName(urlPath.Path, w.bucketName), nil
}

func (*gcpWaitress) hasValidContentType(fileHeader *multipart.FileHeader) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	return acceptedContentTypes[contentType]
}
