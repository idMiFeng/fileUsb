package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Download(c *gin.Context) {
	path := c.Query("path")
	log.Println(path)
	if path == "" {
		c.String(http.StatusBadRequest, "Path parameter is required")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// 设置响应头部信息
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(path))

	// 设置正确的 MIME 类型
	var contentType string
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4":
		contentType = "video/mp4"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".wmv":
		contentType = "video/x-ms-wmv"
	default:
		contentType = "application/octet-stream" // 默认二进制流类型
	}

	c.Header("Content-Type", contentType)

	// 发送文件给客户端
	c.File(path)
}
