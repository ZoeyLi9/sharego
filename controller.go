package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

//定义删除请求的结构体
type DeleteRequest struct {
	Id    int
	Path  string
	Token string //校验符
}

//定义首页以及查询功能
func GetIndex(ctx *gin.Context) {
	query := ctx.Query("search") //路由后定义的关键字

	//定义查询控制功能
	files, _ := Query(query)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"message": "",
		"files":   files,
	})
}

//加载前端静态文件
func GetStaticFile(ctx *gin.Context) {
	path := ctx.Param("file")                          //跟router对应
	ctx.FileFromFS("public/static/"+path, http.FS(fs)) //读取静态文件，打包（编译）静态文件，返回公共路径。
}

//加载第三方样式
func GetLibFile(ctx *gin.Context) {
	path := ctx.Param("file")
	ctx.FileFromFS("public/lib/"+path, http.FS(fs))
}

//上传文件功能控制
func UploadFile(ctx *gin.Context) {
	//将描述信息生成post请求
	description := ctx.PostForm("description")
	if description == "" {
		description = "No Description"
	}

	//将上传者信息生成post请求
	uploader := ctx.PostForm("uploader")
	if uploader == "" {
		uploader = "Unknown User"
	}
	current := time.Now().Format("2006-01-02 21:50:14")

	//使用gin实现多文件上传
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form error: %s", err.Error()))
		return
	}
	//获取参数file后面的多个文件名，存放到数组files中
	files := form.File["file"] //type为“file”
	//获取每个文件的path路径
	for _, file := range files { //不接收key，只接收value
		//获取各文件名
		filename := filepath.Base(file.Filename)
		path := "/upload/" + filename
		//将多个获取到的文件存储到指定文件夹，并进行报错处理
		if err := ctx.SaveUploadedFile(file, UploadPath+"/"+filename); err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file failed: %s", err.Error()))
			return
		}
		//将上述post信息同步更新到数据库中
		fileObject := &File{
			Description: description,
			Uploader:    uploader,
			Time:        current,
			Filename:    filename,
			Path:        path,
		}
		err = fileObject.Insert()
		if err != nil {
			_ = fmt.Errorf(err.Error())
		}
	}
	//链接重定向至首页 mark
	ctx.Redirect(http.StatusSeeOther, "./")
}

//删除文件功能控制
func DeleteFile(ctx *gin.Context) {
	var delreq DeleteRequest
	//对字符串进行解码
	err := json.NewDecoder(ctx.Request.Body).Decode(delreq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Parameter",
		})
		return
	}

	//删除文件和数据库记录
	if *Token == delreq.Token {
		//找到文件对应的id
		fileObject := &File{
			Id: delreq.Id,
		}
		DB.Where("Id = ?", delreq.Id).First(&fileObject)
		//删除文件和记录
		err := fileObject.Delete()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": err.Error(),
			})
		} else {
			var message string
			message = "Record delete successfully"
			if fileObject.IsLocalFile {
				message = "File delete successfully"
			}
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": message,
			})
		}
	} else {
		ctx.JSON(http.StatusOK,gin.H{
			"success": true,
			"message": "Invalid password",
		})
	}
}
