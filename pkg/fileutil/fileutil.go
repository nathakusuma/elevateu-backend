package fileutil

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"mime/multipart"
	"sync"
)

var (
	once             sync.Once
	fileUtilInstance IFileUtil
)

type IFileUtil interface {
	CheckMIMEFileType(file multipart.File, allowed []string) (bool, string, error)
	Upload(ctx context.Context, file io.Reader, path string) (string, error)
	GetSignedURL(path string) (string, error)
	GetUploadSignedURL(path string) (string, error)
	Delete(ctx context.Context, path string) error
}

type fileUtil struct {
	client *storage.Client
}

func NewFileUtil(client *storage.Client) IFileUtil {
	once.Do(func() {
		fileUtilInstance = &fileUtil{
			client: client,
		}
	})
	return fileUtilInstance
}
