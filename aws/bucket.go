package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// BucketService stores images in a private bucket and hands out presigned URLs, which the app can
// render without credentials until they expire. The same object gets the same URL inside a signing
// window of half the expiration, so the app can cache the image by URL.
type BucketService interface {
	// Upload saves the image inside params.Folder and returns its key and presigned URL.
	// The content type is detected from the file content: png, jpeg or webp, up to 4MB.
	Upload(ctx context.Context, params *UploadParams) (*Object, error)

	// GetFile returns the presigned URL of a single image, or ErrNotFound.
	GetFile(ctx context.Context, key string) (*Object, error)

	// ListFolder returns every image directly inside the folder, without its subfolders.
	ListFolder(ctx context.Context, folder string) ([]Object, error)

	// DeleteFile removes a single image. Deleting a key that does not exist is not an error.
	DeleteFile(ctx context.Context, key string) error
}

type s3API interface {
	s3.ListObjectsV2APIClient
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type objectPresigner interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type awsBucket struct {
	s3Client   s3API
	presigner  objectPresigner
	signer     *v4.Signer
	bucketName string
	expiration time.Duration
	now        func() time.Time
}

func NewAwsBucket(configs *BucketConfig) (BucketService, error) {
	if configs.BucketName == "" || configs.Region == "" {
		return nil, errors.New("the bucket name and region are required")
	}

	expiration := configs.URLExpiration
	if expiration == 0 {
		expiration = defaultURLExpiration
	}
	if expiration < 2*time.Second || expiration > maxURLExpiration {
		return nil, fmt.Errorf("invalid URL expiration %s, it must be between 2s and %s", expiration, maxURLExpiration)
	}

	client := s3.New(s3.Options{
		Region: configs.Region,
		Credentials: aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				cred := aws.Credentials{
					AccessKeyID:     configs.AccessKey,
					SecretAccessKey: configs.SecretKey,
				}

				if !cred.HasKeys() {
					return aws.Credentials{}, errors.New("the keys are missing")
				}
				return cred, nil
			},
		),
	})
	return &awsBucket{
		s3Client:  client,
		presigner: s3.NewPresignClient(client),
		// Same options as the S3 default presigner.
		signer: v4.NewSigner(func(so *v4.SignerOptions) {
			so.DisableURIPathEscaping = true
		}),
		bucketName: configs.BucketName,
		expiration: expiration,
		now:        brasiliaTime,
	}, nil
}

func (bucket *awsBucket) Upload(ctx context.Context, params *UploadParams) (*Object, error) {
	if params == nil || params.File == nil || params.FileHeader == nil {
		return nil, errors.New("the file is required")
	}

	defer func(File multipart.File) {
		err := File.Close()
		if err != nil {
			log.Println("[toolkit.Upload] error closing the multipart file")
		}
	}(params.File)

	folder, err := normalizeFolder(params.Folder)
	if err != nil {
		return nil, err
	}

	if params.FileHeader.Size > multipartMaxLength {
		return nil, ErrFileTooLarge
	}

	contentType, err := detectContentType(params.File)
	if err != nil {
		return nil, err
	}

	key := buildKey(folder, contentType)
	_, err = bucket.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket.bucketName),
		Key:           aws.String(key),
		Body:          params.File,
		ContentLength: aws.Int64(params.FileHeader.Size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("error uploading the file to the bucket %s: %w", bucket.bucketName, err)
	}

	object := Object{
		Key:          key,
		Size:         params.FileHeader.Size,
		LastModified: bucket.now(),
	}
	if err = bucket.sign(ctx, &object); err != nil {
		return nil, err
	}
	return &object, nil
}

func (bucket *awsBucket) GetFile(ctx context.Context, key string) (*Object, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	head, err := bucket.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error loading the file %s: %w", key, err)
	}

	object := Object{
		Key:          key,
		Size:         aws.ToInt64(head.ContentLength),
		LastModified: aws.ToTime(head.LastModified),
	}
	if err = bucket.sign(ctx, &object); err != nil {
		return nil, err
	}
	return &object, nil
}

func (bucket *awsBucket) ListFolder(ctx context.Context, folder string) ([]Object, error) {
	folder, err := normalizeFolder(folder)
	if err != nil {
		return nil, err
	}

	prefix := folder + "/"
	paginator := s3.NewListObjectsV2Paginator(bucket.s3Client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket.bucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	objects := make([]Object, 0)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing the folder %s in the bucket %s: %w", folder, bucket.bucketName, err)
		}

		for _, content := range page.Contents {
			key := aws.ToString(content.Key)
			if key == prefix {
				continue // folder marker created by the console
			}

			object := Object{
				Key:          key,
				Size:         aws.ToInt64(content.Size),
				LastModified: aws.ToTime(content.LastModified),
			}
			if err = bucket.sign(ctx, &object); err != nil {
				return nil, err
			}
			objects = append(objects, object)
		}
	}
	return objects, nil
}

func (bucket *awsBucket) DeleteFile(ctx context.Context, key string) error {
	if err := validateKey(key); err != nil {
		return err
	}

	_, err := bucket.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error deleting the file %s: %w", key, err)
	}
	return nil
}

// sign fills the presigned URL of the object, signed at the start of the current signing window.
func (bucket *awsBucket) sign(ctx context.Context, object *Object) error {
	signedAt := bucket.now().Truncate(signingWindow(bucket.expiration))

	request, err := bucket.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket.bucketName),
		Key:    aws.String(object.Key),
	}, func(options *s3.PresignOptions) {
		options.Expires = bucket.expiration
		options.Presigner = fixedTimePresigner{signer: bucket.signer, signedAt: signedAt.UTC()}
	})
	if err != nil {
		return fmt.Errorf("error signing the URL of the file %s: %w", object.Key, err)
	}

	object.URL = request.URL
	object.ExpiresAt = signedAt.Add(bucket.expiration)
	return nil
}
