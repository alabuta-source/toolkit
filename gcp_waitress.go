package toolkit

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	defaultCacheControl = "Cache-Control:no-cache, max-age=0"
)

type GCPWaitressManager interface {
	// SaveFile saves a file to the bucket and returns the name of the file or an error
	// The prefix is used to create a folder in the bucket, if the prefix is empty, the file will be saved in the root of the bucket.
	// Use the / after the prefix to create a folder, example: "images/"
	// The cacheControl is used to set the cache control header in the file, if the cacheControl is empty, the default value will be used.
	// The default value is: "Cache-Control:no-cache, max-age=0"
	SaveFile(file multipart.File, fileHeader *multipart.FileHeader, prefix, cacheControl string) (string, error)

	// ListFiles returns a list of files in the bucket
	// The prefix is used to filter the files, if the prefix is empty, all files will be returned.
	// Make sure to add the / after the prefix to filter the files, example: "images/"
	ListFiles(prefix string) ([]string, error)

	// DeleteFile use this method to delete a single file.
	DeleteFile(filename string) error
}

type gcpWaitress struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

// NewGCPWaitress creates a new GCP Waitress pointer to access the GCP Bucket
// The bucket is the name of the bucket to access in Google Cloud Storage
// The request is the http request to get the context, this is used to create the context for the GCP Client
// The gcpKey is the json key to access the bucket, this key can be created in the Google Cloud Console and should have the Storage Admin role
func NewGCPWaitress(bucket string, request *http.Request, gcpKey []byte) (GCPWaitressManager, error) {
	ctx := appengine.NewContext(request)
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(gcpKey))
	if err != nil {
		return nil, err
	}

	return &gcpWaitress{
		client: client,
		bucket: client.Bucket(bucket),
		ctx:    ctx,
	}, nil
}

func (w *gcpWaitress) SaveFile(file multipart.File, fileHeader *multipart.FileHeader, prefix, cacheControl string) (string, error) {
	obj := w.bucket.
		Object(fileHeader.Filename).
		NewWriter(w.ctx)

	if cacheControl == "" {
		cacheControl = defaultCacheControl
	}
	obj.Attrs().CacheControl = cacheControl
	obj.Attrs().Prefix = prefix

	if _, err := io.Copy(obj, file); err != nil {
		return "", fmt.Errorf("unable to write file to Google Storage: %w", err)
	}

	if er := obj.Close(); er != nil {
		return "", fmt.Errorf("unable to close the writer: %w", er)
	}
	return obj.Attrs().Name, nil
}

func (w *gcpWaitress) ListFiles(prefix string) ([]string, error) {
	resp := make([]string, 0)
	objects := w.bucket.Objects(w.ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: "/",
	})

	for {
		attrs, err := objects.Next()
		if err != nil {
			return nil, fmt.Errorf("error loading file: %w", err)
		}
		if attrs == nil {
			break
		}
		resp = append(resp, attrs.Name)
	}
	return resp, nil
}

func (w *gcpWaitress) DeleteFile(filename string) error {
	oHandle := w.bucket.Object(filename)

	attrs, err := oHandle.Attrs(w.ctx)
	if err != nil {
		return fmt.Errorf("error getting meta information about the object to delete it: %w", err)
	}

	if er := oHandle.
		If(storage.Conditions{GenerationMatch: attrs.Generation}).
		Delete(w.ctx); er != nil {
		return fmt.Errorf("error trying delete obj: %w", er)
	}
	return nil
}
