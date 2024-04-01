package controller

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func copyFiles(src string, destDir string) error {
	// 打开源文件或目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 根据源文件或目录的类型执行相应的复制操作
	if srcInfo.IsDir() {
		// 复制目录
		return copyDirectory(src, destDir)
	} else {
		// 复制文件
		return copyFile(src, destDir)
	}
}

func copyFile(srcFile string, destDir string) error {
	// 打开源文件
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	destFile := filepath.Join(destDir, filepath.Base(srcFile))
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	// 复制文件内容
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	log.Println("File copied successfully:", destFile)
	return nil
}

func copyDirectory(srcDir string, destDir string) error {
	// 创建目标目录
	dest := filepath.Join(destDir, filepath.Base(srcDir))
	err := os.Mkdir(dest, os.ModePerm)
	if err != nil {
		return err
	}

	// 获取源目录下的所有文件和子目录
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// 递归复制文件和子目录
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			err = copyDirectory(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// 复制文件
			err = copyFile(srcPath, dest)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Directory copied successfully:", dest)
	return nil
}
