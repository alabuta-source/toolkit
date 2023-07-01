package toolkit

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"io"
	"mime/multipart"
	"net/http"
)

type GCPWaitressManager interface {
	SaveFile(multipart.File, *multipart.FileHeader) (string, error)
	ListFiles(bucket string) ([]string, error)
}

type gcpWaitress struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

func NewGCPWaitress(bucket string, request *http.Request, gcpKey *GCPBucketAuthJson) (GCPWaitressManager, error) {
	ctx := appengine.NewContext(request)
	bytes, jErr := json.Marshal(gcpKey)
	if jErr != nil {
		return nil, fmt.Errorf("error trying marchal gcp key :%s", jErr.Error())
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(bytes))
	if err != nil {
		return nil, err
	}

	return &gcpWaitress{
		client: client,
		bucket: client.Bucket(bucket),
		ctx:    ctx,
	}, nil
}

func (w *gcpWaitress) SaveFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	obj := w.bucket.
		Object(fileHeader.Filename).
		NewWriter(w.ctx)

	obj.Attrs().CacheControl = "Cache-Control:no-cache, max-age=0"

	if _, err := io.Copy(obj, file); err != nil {
		return "", fmt.Errorf("unable to write file to Google Storage: %s", err.Error())
	}

	if er := obj.Close(); er != nil {
		return "", fmt.Errorf("unable to close the writer: %s", er.Error())
	}
	return obj.Attrs().Name, nil
}

func (w *gcpWaitress) ListFiles(bucket string) ([]string, error) {
	resp := make([]string, 0)
	objects := w.client.Bucket(bucket).Objects(w.ctx, nil)

	for {
		attrs, err := objects.Next()
		if attrs == nil {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error loading file: %s", err.Error())
		}
		resp = append(resp, attrs.Name)
	}
	return resp, nil
}
