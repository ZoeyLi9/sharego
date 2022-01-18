package main

import "github.com/gin-gonic/gin"

func SetIndexRouter(router *gin.Engine) { //封装好的路由引擎，用于定义服务路由信息、组装插件、运行服务。
	//上传文件的路径
	router.Static("/upload", "./upload")
	//首页路径
	router.GET("/", GetIndex) //http get
	//静态文件路径
	router.GET("/public/static/:file", GetStaticFile)
	//第三方样式库文件路径
	router.GET("/public/lib/:file", GetLibFile)
}

//定义Upload和Delete接口路径
func SetApiRouter(router *gin.Engine) {
	router.POST("/upload", UploadFile)
	router.POST("/delete", DeleteFile)
}