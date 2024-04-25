package controller

import (
	"fileUsb/service"
	"fileUsb/utils"
	"github.com/shirou/gopsutil/disk"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// DiskInfo 获取U盘信息
func DiskInfo() []utils.Blockdevice {
	Blockdevices := make([]utils.Blockdevice, 0)
	Blockdevices = utils.ReadForensicDisk()
	return Blockdevices
}

// ListDisk 获取U盘的目录和文件
func ListDisk(rootPath string) ([]service.FileInfo, error) {
	var fileList []service.FileInfo

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileInfo := service.FileInfo{
			Name:     file.Name(),
			Path:     filepath.Join(rootPath, file.Name()),
			IsDir:    file.IsDir(),
			Modified: file.ModTime().Unix(),
			Size:     file.Size(),
		}
		fileList = append(fileList, fileInfo)

	}

	return fileList, nil
}

// ListMountpoint 获取所有U盘的挂载点路径
func ListMountpoint() []string {
	var Mountpoints []string
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	for _, partition := range partitions {
		if strings.Contains(partition.Device, "sdb") {
			Mountpoints = append(Mountpoints, partition.Mountpoint)
		}
	}
	return Mountpoints
}

func SearchFiles(rootPath string, searchQuery string) ([]service.FileInfo, error) {
	var fileList []service.FileInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if containsAnyCase(info.Name(), searchQuery) {
				fileInfo := service.FileInfo{
					Name:     info.Name(),
					Path:     path,
					IsDir:    info.IsDir(),
					Modified: info.ModTime().Unix(),
					Size:     info.Size(),
				}
				fileList = append(fileList, fileInfo)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

// containsAnyCase 检查字符串s中是否包含子字符串substr，忽略大小写。
func containsAnyCase(s, substr string) bool {
	// 先尝试使用快速的英文字符匹配
	if strings.Contains(strings.ToLower(s), strings.ToLower(substr)) {
		return true
	}

	// 将子字符串转换为小写
	substrLower := strings.ToLower(substr)

	// 遍历s中的每个字符
	for len(s) > 0 {
		// 解析第一个字符
		r, size := utf8.DecodeRuneInString(s)
		// 移动指针
		s = s[size:]

		// 如果是英文字符，跳过
		if r < utf8.RuneSelf && unicode.IsLetter(r) {
			continue
		}

		// 如果是中文字符，直接比较
		if unicode.Is(unicode.Han, r) {
			// 如果中文字符等于给定的中文字符，返回true
			if strings.ContainsRune(substrLower, r) {
				return true
			}
		}
	}

	return false
}
