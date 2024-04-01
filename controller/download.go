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
	c.Header("Content-Type", "application/octet-stream") // 通用的二进制流类型
	// 如果文件类型是视频，可以根据实际情况设置 MIME 类型
	if strings.HasSuffix(strings.ToLower(path), ".mp4") {
		c.Header("Content-Type", "video/mp4")
	}

	// 发送文件给客户端
	c.File(path)
}
