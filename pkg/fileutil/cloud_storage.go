package fileutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
)

func (u *fileUtil) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	bucket := env.GetEnv().GCPStorageBucketName

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := u.client.Bucket(bucket).Object(path).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, path)

	return url, nil
}

func (u *fileUtil) GetFullURL(path string) string {
	bucket := env.GetEnv().GCPStorageBucketName
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, path)
}

func (u *fileUtil) GetSignedURL(path string) (string, error) {
	bucket := env.GetEnv().GCPStorageBucketName

	url, err := u.client.Bucket(bucket).SignedURL(path, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		return "", fmt.Errorf("client.Bucket(%q).SignedURL: %w", bucket, err)
	}

	return url, nil
}

func (u *fileUtil) GetUploadSignedURL(path, contentType string) (string, error) {
	bucket := env.GetEnv().GCPStorageBucketName

	options := &storage.SignedURLOptions{
		Method:  http.MethodPut,
		Expires: time.Now().Add(10 * time.Minute),
	}

	if contentType != "" {
		options.ContentType = contentType
	}

	url, err := u.client.Bucket(bucket).SignedURL(path, options)
	if err != nil {
		return "", fmt.Errorf("client.Bucket(%q).SignedURL: %w", bucket, err)
	}

	return url, nil
}

func (u *fileUtil) Delete(ctx context.Context, path string) error {
	bucket := env.GetEnv().GCPStorageBucketName

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := u.client.Bucket(bucket).Object(path)

	// Set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err2 := o.Delete(ctx); err2 != nil {
		return fmt.Errorf("Object(%q).Delete: %w", path, err2)
	}

	return nil
}
