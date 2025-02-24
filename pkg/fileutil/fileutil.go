package fileutil

import (
	"mime/multipart"
	"sync"
)

var (
	once             sync.Once
	fileUtilInstance IFileUtil
)

type IFileUtil interface {
	CheckMIMEFileType(file multipart.File, allowed []string) (bool, string, error)
}

type fileUtil struct{}

func NewFileUtil() IFileUtil {
	once.Do(func() {
		fileUtilInstance = &fileUtil{}
	})
	return fileUtilInstance
}
