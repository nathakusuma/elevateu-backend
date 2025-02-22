package fileutil

import (
	"net/http"
	"time"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/internal/infra/gcp"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

func GetSignedURL(path string) string {
	bucket := env.GetEnv().GCPStorageBucketName

	client := gcp.Client
	if client == nil {
		log.Error(nil, "[GCP][GetSignedURL] GCP client is nil")
		return ""
	}

	url, err := client.Bucket(bucket).SignedURL(path, &storage.SignedURLOptions{
		Method:  http.MethodGet,
		Expires: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		log.Error(map[string]interface{}{
			"error": err,
		}, "[GCP][GetSignedURL] Failed to get signed url")
		return ""
	}

	return url
}
