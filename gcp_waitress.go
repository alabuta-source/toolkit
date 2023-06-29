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
	"net/url"
	"path"
)

type GCPWaitressManager interface {
	SaveFile(multipart.File, *multipart.FileHeader) (string, error)
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
	obj := w.bucket.Object(fileHeader.Filename)
	wc := obj.NewWriter(w.ctx)

	wc.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0"

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("unable to write file to Google Storage: %s", err.Error())
	}

	if er := wc.Close(); er != nil {
		return "", fmt.Errorf("unable to close the writer: %s", er.Error())
	}

	return obj.ObjectName(), nil
}

func objectNameFromUrl(imgUrl string) (string, error) {
	if imgUrl == "" {
		return "", nil
	}

	urlPath, err := url.Parse(imgUrl)
	if err != nil {
		return "", fmt.Errorf("unable to parse the url: %s", err.Error())
	}
	return path.Base(urlPath.Path), nil
}
