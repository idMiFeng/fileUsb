package controller

import (
	"fileUsb/utils"
	"github.com/shirou/gopsutil/disk"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// DiskInfo 获取U盘信息
func DiskInfo() []utils.Blockdevice {
	Blockdevices := make([]utils.Blockdevice, 0)
	Blockdevices = utils.ReadForensicDisk()
	return Blockdevices
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

// ListDisk 获取U盘的目录和文件
func ListDisk(rootPath string) ([]FileInfo, error) {
	var fileList []FileInfo

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileInfo := FileInfo{
			Name:  file.Name(),
			Path:  filepath.Join(rootPath, file.Name()),
			IsDir: file.IsDir(),
		}
		fileList = append(fileList, fileInfo)

		// 如果是目录，则递归获取目录下的文件信息
		if file.IsDir() {
			subList, err := ListDisk(filepath.Join(rootPath, file.Name()))
			if err != nil {
				log.Printf("Error listing directory %s: %v\n", fileInfo.Path, err)
				continue
			}
			fileList = append(fileList, subList...)
		}
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
