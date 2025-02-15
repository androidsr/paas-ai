package toolkit

import (
	"net"
	"strings"
	"time"
)

func Signature(macAddress string, publish string) bool {
	for _, s := range strings.Split(macAddress, ",") {
		if strings.Contains(getMacAddress(), s) {
			return true
		}
	}
	return isWithin30Days(publish)
}

// 获取本机网卡地址
func getMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	// 遍历网卡列表，查找第一个有效的网卡
	var hardwareAddr string
	for _, iface := range interfaces {
		hardwareAddr += strings.ReplaceAll(strings.ToUpper(iface.HardwareAddr.String()), ":", "-")
	}
	return hardwareAddr
}

// 比较时间是否在指定时间之后的30天内
func isWithin30Days(specifiedTime string) bool {
	layout := "2006-01-02"
	parsedTime, err := time.ParseInLocation(layout, specifiedTime, time.Local) // 使用本地时区
	if err != nil {
		return false
	}

	currentTime := time.Now().In(time.Local) // 获取本地时间（中国时区）
	return currentTime.Sub(parsedTime) <= 30*24*time.Hour
}
