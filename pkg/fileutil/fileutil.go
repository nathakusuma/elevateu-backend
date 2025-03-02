package fileutil

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"sync"

	"cloud.google.com/go/storage"

	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

var (
	once             sync.Once
	fileUtilInstance IFileUtil
)

type IFileUtil interface {
	CheckMIMEFileType(file multipart.File, allowed []string) (bool, string, error)
	Upload(ctx context.Context, file io.Reader, path string) (string, error)
	GetFullURL(path string) string
	GetSignedURL(path string) (string, error)
	GetUploadSignedURL(path, contentType string) (string, error)
	Delete(ctx context.Context, path string) error
	ValidateAndUploadFile(ctx context.Context, header *multipart.FileHeader, allowedTypes []string,
		path string) (string, error)
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

func (u *fileUtil) ValidateAndUploadFile(ctx context.Context, header *multipart.FileHeader, allowedTypes []string,
	path string) (string, error) {
	if header.Size > 2*MegaByte {
		return "", errorpkg.ErrFileTooLarge.WithDetail(
			fmt.Sprintf("File size is too large (%s). Please upload a file less than 2MB",
				ByteToAppropriateUnit(header.Size)))
	}

	file, err := header.Open()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"path":  path,
		}, "[FileUtil][ValidateAndUploadFile] Failed to open file")
		return "", errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	defer file.Close()

	ok, fileType, err := u.CheckMIMEFileType(file, allowedTypes)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"path":  path,
		}, "[FileUtil][ValidateAndUploadFile] Failed to check MIME file type")
		return "", errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	if !ok {
		return "", errorpkg.ErrInvalidFileFormat.WithDetail(
			fmt.Sprintf("File type %s is not allowed. Please upload a valid file", fileType))
	}

	url, err := u.Upload(ctx, file, path)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"path":  path,
		}, "[FileUtil][ValidateAndUploadFile] Failed to upload file")
		return "", errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return url, nil
}
