package contract

import (
	"context"
	"io"
)

type IStorageRepository interface {
	Upload(ctx context.Context, file io.Reader, path string) (string, error)
	Delete(ctx context.Context, path string) error
}
