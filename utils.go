package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

//打开浏览器
func openBrowser(url string) {
	var err error
	//判断用户的os类型
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Println(err)
	}
}

//获取ip地址
func getIp() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return
	}

	for _, a := range addrs {
		//判断该ip地址是否为回环地址
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil { //说明该ip地址是一个ipv4地址
				ip = ipNet.IP.String()
				//判断该ip地址是否为本机地址（即以192.168开始），且不是网关保留地址（不能结束为.1），则返回该ip地址，否则返回空
				if strings.HasPrefix(ip, "192.168") && !strings.HasSuffix(ip, ".1") {
					return
				}
				ip = ""
			}
		}
	}
	return
}

var sizeKB = 1024
var sizeMB = sizeKB * 1024
var sizeGB = sizeMB * 1024
var sizeTB = sizeGB * 1024

//将文件大小字节转换为KB、MB形式
func Bytes2Size(num int64) string {
	numS := ""
	unit := "B"
	if num/int64(sizeTB) >= 1 {
		numS = fmt.Sprintf("%f", float64(num)/float64(sizeTB))
		unit = "TB"
	} else if num/int64(sizeGB) >= 1 {
		numS = fmt.Sprintf("%f", float64(num)/float64(sizeGB))
		unit = "GB"
	} else if num/int64(sizeMB) >= 1 {
		numS = fmt.Sprintf("%f", float64(num)/float64(sizeMB))
		unit = "MB"
	} else if num/int64(sizeKB) >= 1 {
		numS = fmt.Sprintf("%f", float64(num)/float64(sizeKB))
		unit = "KB"
	} else {
		numS = fmt.Sprintf("%d", num)
	}
	numS = strings.Split(numS, ".")[0]
	return numS + " " + unit
}
