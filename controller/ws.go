package controller

import (
	"encoding/json"
	"fileUsb/error"
	"fileUsb/service"
	"github.com/olahol/melody"
	"log"
)

type WsRes struct {
	Type       string   `json:"type"`
	Mountpoint string   `json:"mountpoint"`
	Path       []string `json:"path"`
	SortType   string   `json:"sort_type"`
	SearchName string   `json:"search_name"`
}

type WsMessage struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Op   string      `json:"op,omitempty"`
}

const (
	successCode = 2000
	errorCode   = 4000
)

func HandleWebSocketMessage(s *melody.Session, msg []byte) {
	log.Println("Received message:", string(msg))
	// 前端发送的数据格式为 JSON：{"type": "download", "path": "/media/fengmi/3A5EBF185EBECC3F"}
	var data WsRes
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("Error decoding JSON:", err)
		ws := WsMessage{
			Code: errorCode,
			Msg:  err.Error(),
		}
		data, _ := json.Marshal(&ws)
		_ = s.Write(data)
		return
	}

	// 根据指令类型执行操作
	switch data.Type {
	case "search":
		files, _ := SearchFiles(data.Mountpoint, data.SearchName)
		data, _ := json.Marshal(&files)
		ws := WsMessage{
			Code: successCode,
			Data: files,
			Op:   "search",
		}
		data, _ = json.Marshal(&ws)
		_ = s.Write(data)
	case "size", "name", "data", "file_type":
		files, err := ListDisk(data.Mountpoint)
		if err != nil {
			log.Println("Error listing disk:", err)
			ws := WsMessage{
				Code: errorCode,
				Msg:  error.ErrFileFalse.Error(),
			}
			data, _ := json.Marshal(&ws)
			_ = s.Write(data)
			return
		}
		// 对文件列表进行排序
		sortedFiles, _ := service.SortFiles(data.Type, data.SortType, files)
		ws := WsMessage{
			Code: successCode,
			Data: sortedFiles,
			Op:   "sort",
		}
		data, _ := json.Marshal(&ws)
		_ = s.Write(data)
	case "disk":
		Blockdevice := DiskInfo()
		ws := WsMessage{
			Code: successCode,
			Data: Blockdevice,
			Op:   "disk",
		}
		data, _ := json.Marshal(&ws)
		_ = s.Write(data)

	case "info":
		//log.Println(data.Mountpoint)
		if datas, err := ListDisk(data.Mountpoint); err != nil {
			ws := WsMessage{
				Code: errorCode,
				Msg:  err.Error(),
			}
			data, _ := json.Marshal(&ws)
			_ = s.Write(data)
		} else {
			ws := WsMessage{
				Code: successCode,
				Data: datas,
				Op:   "info",
			}
			data, _ := json.Marshal(&ws)
			_ = s.Write(data)
		}

	case "copy":
		// /mnt4/chroot/device_data/usb_data/2024-03-29/15-30-20
		// usb_data    插件名称
		// 2024-03-29   日期
		// 15-30-20 时间
		destDir := "files"
		for _, path := range data.Path {
			if err := copyFiles(path, destDir); err != nil {
				ws := WsMessage{
					Code: errorCode,
					Msg:  err.Error(),
				}
				data, _ := json.Marshal(&ws)
				_ = s.Write(data)
			} else {
				log.Println("copy success", path)
			}
		}
		ws := WsMessage{
			Code: successCode,
			Data: "success",
			Op:   "copy",
		}
		data, _ := json.Marshal(&ws)
		_ = s.Write(data)
	default:
		log.Println("Unknown command")
	}
}

func WsHandleConnect(s *melody.Session) {
	ws := WsMessage{
		Code: successCode,
		Msg:  "",
	}
	data, _ := json.Marshal(&ws)
	_ = s.Write(data)
}

func WsHandleDisconnect(s *melody.Session) {
	ws := WsMessage{
		Code: successCode,
		Msg:  "",
	}
	data, _ := json.Marshal(&ws)
	_ = s.Write(data)
}
