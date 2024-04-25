package service

import (
	"fileUsb/error"
	"log"
	"sort"
	"strings"
)

type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Modified int64  `json:"modified"` // 修改日期
	Size     int64  `json:"size"`     // 文件大小
}

const (
	// SortByName 表示按名称排序
	SortByName     = "name"
	SortBySize     = "size"
	SortByData     = "data"
	SortByFileTyle = "file_type"
	Asc            = "asc"
	Desc           = "desc"
)

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}
func SortFiles(Type string, sortType string, files []FileInfo) ([]FileInfo, interface{}) {
	switch {
	case sortType == Desc && Type == SortByName:
		sort.SliceStable(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name) > strings.ToLower(files[j].Name)
		})
	case sortType == Asc && Type == SortByName:
		sort.SliceStable(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
		})
	case sortType == Asc && Type == SortBySize:
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Size < files[j].Size
		})
	case sortType == Desc && Type == SortBySize:
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Size > files[j].Size
		})
	case sortType == Asc && Type == SortByData:
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Modified < files[j].Modified
		})
	case sortType == Desc && Type == SortByData:
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Modified > files[j].Modified
		})
	case sortType == Desc && Type == SortByFileTyle:
		sort.SliceStable(files, func(i, j int) bool {
			return getFileExtension(files[i].Name) < getFileExtension(files[j].Name)
		})
	case sortType == Asc && Type == SortByFileTyle:
		sort.SliceStable(files, func(i, j int) bool {
			return getFileExtension(files[i].Name) > getFileExtension(files[j].Name)
		})
	default:
		log.Println("无效的排序参数", Type, sortType)
		return nil, error.ErrUnauthorized.Error() // 处理未匹配到任何条件的情况
	}

	return files, nil

}
