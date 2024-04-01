package main

import (
	"fileUsb/controller"
	"fileUsb/global"
	"fileUsb/router"
	"github.com/olahol/melody"
	"log"
	"time"
)

func main() {
	OtherInit()
	router.InitRouter()
}

func OtherInit() {
	global.M = melody.New()
	global.M.Config = &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          60 * time.Second,
		PingPeriod:        (60 * time.Second * 9) / 10,
		MaxMessageSize:    51200,
		MessageBufferSize: 256,
	}
	global.M.HandleConnect(controller.WsHandleConnect)
	global.M.HandleDisconnect(controller.WsHandleDisconnect)
	global.M.HandleMessage(controller.HandleWebSocketMessage)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
}
