package main

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
	"strconv"
)

var (
	host  = flag.String("host", "localhost", "please input the server ip address")
	port  = flag.Int("port", 8080, "please input the server port")
	Token = flag.String("token", "123456", "Set private password")
	path  = flag.String("path", "", "Set the folder to be shared")
)

var Url = ""
var UploadPath = "./upload"
var LocalRoot = UploadPath

//go:embed public
var fs embed.FS

//若upload文件夹不存在则新建文件夹
func init() {
	if _, err := os.Stat(UploadPath); os.IsNotExist(err) {
		_ = os.Mkdir(UploadPath, 0777)
	}
}

//加载HTML模板文件
func loadTemplate() *template.Template {
	t := template.Must(template.New("").ParseFS(fs,"public/*.html"))
	return t
}

func main() {
	flag.Parse()
	//设置gin为发行模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	//初始化数据库
	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := gin.Default()
	server.SetHTMLTemplate(loadTemplate())

	SetApiRouter(server)
	SetIndexRouter(server)

	//判断路径
	if *path != "" { //若用户已自定义上传路径
		LocalRoot = *path
	}

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

	err = server.Run(":" + realPort)
	if err != nil {
		log.Println(err)
	}
}
