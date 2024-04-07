package router

import (
	"fileUsb/controller"
	"fileUsb/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "404 ,page not exists!",
	})
}
func InitRouter() {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.NoRoute(NotFound)
	v1 := r.Group("api")
	{

		//v1.GET("/disk_info", controller.DiskInfo)
		v1.GET("/download", controller.Download)
		v1.GET("/ws", controller.Copy)

	}
	_ = r.Run(":3016")
}
