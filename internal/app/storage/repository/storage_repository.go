package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
)

type storageRepository struct {
	client *storage.Client
}

func NewStorageRepository(client *storage.Client) contract.IStorageRepository {
	return &storageRepository{
		client: client,
	}
}

func (r *storageRepository) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	bucket := env.GetEnv().GCPStorageBucketName

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := r.client.Bucket(bucket).Object(path).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, path)

	return url, nil
}

func (r *storageRepository) GetSignedURL(originalURL string) (string, error) {
	bucket := env.GetEnv().GCPStorageBucketName
	expectedPrefix := "https://storage.googleapis.com/" + bucket + "/"

	// Validate URL format
	if !strings.HasPrefix(originalURL, expectedPrefix) {
		return "", fmt.Errorf("invalid URL format: must start with %s", expectedPrefix)
	}

	// Extract and validate path
	path := originalURL[len(expectedPrefix):]
	if path == "" {
		return "", errors.New("empty path is not allowed")
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return "", errors.New("path traversal attempts are not allowed")
	}

	// Validate path characters
	if !fileutil.IsValidPath(path) {
		return "", errors.New("path contains invalid characters")
	}

	url, err := r.client.Bucket(bucket).SignedURL(path, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		return "", fmt.Errorf("client.Bucket(%q).SignedURL: %w", bucket, err)
	}

	return url, nil
}

func (r *storageRepository) Delete(ctx context.Context, path string) error {
	bucket := env.GetEnv().GCPStorageBucketName

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := r.client.Bucket(bucket).Object(path)

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
