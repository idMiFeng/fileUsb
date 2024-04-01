package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Blockdevices struct {
	Blockdevices []Blockdevice `json:"blockdevices"`
}
type Partition struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Fsavail    string `json:"fsavail"`
	Fssize     string `json:"fssize"`
	Fstype     string `json:"fstype"`
	Fsused     string `json:"fsused"`
	Fsuse      string `json:"fsuse%"`
	Mountpoint string `json:"mountpoint"`
	Label      string `json:"label"`
}
type Blockdevice struct {
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Fsavail    string      `json:"fsavail"`
	Fssize     string      `json:"fssize"`
	Fsused     string      `json:"fsused"`
	Fstype     string      `json:"fstype"`
	Fsuse      string      `json:"fsuse%"`
	Mountpoint string      `json:"mountpoint"`
	Label      string      `json:"label"`
	Partition  []Partition `json:"children"`
}

func matchMountPoint(point string) bool {
	return strings.HasPrefix(point, "/media/")
}

func dealDevice(device Blockdevice, n int) Blockdevice {
	var resDevice Blockdevice
	resDevice = device

	resDevice.Name = "磁盘" + strconv.Itoa(n)
	partition := Partition{
		Name:       "分区0",
		Path:       resDevice.Path,
		Fsavail:    resDevice.Fsavail,
		Fssize:     resDevice.Fssize,
		Fstype:     resDevice.Fstype,
		Fsused:     resDevice.Fsused,
		Fsuse:      resDevice.Fsuse,
		Mountpoint: resDevice.Mountpoint,
		Label:      resDevice.Label,
	}
	resDevice.Partition = []Partition{partition}
	return resDevice

}

func dealPartition(device Blockdevice, n int) Blockdevice {
	var resDevice Blockdevice
	resDevice = device

	var fssize int
	var fsused int
	var fsavail int

	for _, partition := range device.Partition {
		partitionFssize, _ := strconv.Atoi(partition.Fssize)
		partitionFsused, _ := strconv.Atoi(partition.Fsused)
		partitionFsavail, _ := strconv.Atoi(partition.Fsavail)
		fssize += partitionFssize
		fsused += partitionFsused
		fsavail += partitionFsavail

	}
	resDevice.Name = "磁盘" + strconv.Itoa(n)
	resDevice.Fssize = strconv.Itoa(fssize)
	resDevice.Fsused = strconv.Itoa(fsused)
	resDevice.Fsavail = strconv.Itoa(fsavail)
	for m, partition := range resDevice.Partition {
		if partition.Label == "" {
			resDevice.Partition[m].Name = "分区" + strconv.Itoa(m)
		} else {
			resDevice.Partition[m].Name = partition.Label
		}
	}
	return resDevice

}

func ReadForensicDisk() []Blockdevice {
	var blockdeviceList Blockdevices
	var returnDeviceList []Blockdevice
	outContext, _ := RunCommandWithOutput(0, "",
		"bash", "-c", "'lsblk -f -J -e 7,11,3,2 -o NAME,PATH,FSAVAIL,FSSIZE,FSTYPE,FSUSED,FSUSE%,MOUNTPOINT,LABEL -b'")
	err := json.Unmarshal(outContext, &blockdeviceList)
	if err != nil {
		fmt.Println("JSON 解码出错:", err)
		return nil
	}
	for n, device := range blockdeviceList.Blockdevices {
		if matchMountPoint(device.Mountpoint) {
			foramtDevice := dealDevice(device, n)
			returnDeviceList = append(returnDeviceList, foramtDevice)
			continue
		}
		for _, partition := range device.Partition {
			if !matchMountPoint(partition.Mountpoint) {
				continue
			}
			foramtDevice := dealPartition(device, n)
			returnDeviceList = append(returnDeviceList, foramtDevice)
			break
		}
	}
	return returnDeviceList
}

type WsDiskInfoMsg struct {
	Code    int           `json:"code"`
	Operate string        `json:"operate"`
	Msg     string        `json:"msg"`
	Data    []Blockdevice `json:"data"`
}

func ReadStorageDevices(conn *ServerConn, operateType string) {
	for !conn.Flag {
		blockDevices := ReadForensicDisk()
		if len(blockDevices) > 0 {
			wsDeviceInfoMsg := WsDiskInfoMsg{
				Code:    200,
				Operate: operateType,
				Msg:     "success",
				Data:    blockDevices,
			}
			successInfo, _ := json.Marshal(wsDeviceInfoMsg)
			conn.Send <- successInfo
		}
		time.Sleep(2 * time.Second)
	}
}
