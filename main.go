package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	host = flag.String("host", "localhost", "please input the server ip address")
	port = flag.Int("port", 8080, "please input the server port")
)

var Url = " "

func main() {

	//获取ip地址
	if *host == "localhost" {
		ip := getIp()
		if ip != "" {
			*host = ip
		}
	}

	//获取端口号
	var realPort = os.Getenv("PORT")
	if realPort == "" {
		realPort = strconv.Itoa(*port)
	}

	Url = "http://" + *host + ":" + realPort + "/"
	//打开浏览器，进入程序
	openBrowser(Url)

}
