package gcp

import (
	"context"
	"sync"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

var (
	once   sync.Once
	client *storage.Client
)

func NewStorageClient() *storage.Client {
	once.Do(func() {
		ctx := context.Background()
		cl, err := storage.NewGRPCClient(ctx)
		if err != nil {
			log.Fatal(map[string]interface{}{
				"error": err,
			}, "[GCP][NewStorageClient] Failed to create new cloud storage client")
			return
		}

		client = cl
	})

	return client
}
