package controller

import (
	"fileUsb/global"
	"github.com/gin-gonic/gin"
	"log"
)

func Copy(c *gin.Context) {
	if err := global.M.HandleRequest(c.Writer, c.Request); err != nil {
		log.Println("升级到ws协议出现错误:", err)
		return
	}
}
