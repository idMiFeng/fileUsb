package error

import "errors"

var (
	ErrFileFalse    = errors.New("文件损坏")
	ErrUnauthorized = errors.New("文件未排序")
)
