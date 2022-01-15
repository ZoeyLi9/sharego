package main

import (
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
