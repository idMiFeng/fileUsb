package utils

import (
	"golang.org/x/net/websocket"
	"sync"
)

type ServerConn struct {
	WSMutex     sync.Mutex      // WS写锁
	WS          *websocket.Conn // websocket连接
	Exit        chan bool       //退出
	Flag        bool            //判断websocket是否已经关闭
	Receive     chan []byte     //接收客户端的数据
	Send        chan []byte     //发送数据
	MessageType int
}
