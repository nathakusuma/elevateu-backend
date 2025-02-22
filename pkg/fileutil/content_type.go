package fileutil

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

var ImageContentTypes = []string{
	"image/apng",
	"image/avif",
	"image/bmp",
	"image/gif",
	"image/vnd.microsoft.icon",
	"image/jpeg",
	"image/png",
	"image/svg+xml",
	"image/tiff",
	"image/webp",
}

func CheckMIMEFileType(file multipart.File, allowed []string) (bool, string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, "", err
	}

	// Reset file position to beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, "", fmt.Errorf("failed to reset file position: %w", err)
	}

	fileType := http.DetectContentType(buffer)
	for _, allowedType := range allowed {
		if fileType == allowedType {
			return true, fileType, nil
		}
	}

	return false, fileType, nil
}
