package gcp

import (
	"context"
	"sync"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

var (
	once   sync.Once
	Client *storage.Client
)

func NewStorageClient() *storage.Client {
	once.Do(func() {
		ctx := context.Background()
		cl, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatal(map[string]interface{}{
				"error": err,
			}, "[GCP][NewStorageClient] Failed to create new cloud storage client")
			return
		}

		Client = cl
	})

	return Client
}
