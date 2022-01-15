package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIndex(ctx *gin.Context) {
	query := ctx.Query("search") //路由后定义的关键字

	//定义查询控制功能
	files, _ := Query(query)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"message": "",
		"files": files,
	})
}
