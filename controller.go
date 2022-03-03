package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	query := ctx.Query("search") //路由后定义的关键字。获取搜索内容

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
	current := time.Now().Format("2006/01/02 15:04:05")

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
			Link:        path,
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
	err := json.NewDecoder(ctx.Request.Body).Decode(&delreq)
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
		DB.Where("id = ?", delreq.Id).First(&fileObject)
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
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Invalid password",
		})
	}
}

//浏览文件功能控制
func GetExplorerFile(ctx *gin.Context) {
	var localFiles []LocalFile //文件夹数组
	var localf []LocalFile //文件数组

	path := ctx.DefaultQuery("path", "/")
	path, _ = url.QueryUnescape(path) //还原转码后的字符串

	//获取根路径
	rootpath := filepath.Join(LocalRoot, path)
	root, err := os.Stat(rootpath)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": err.Error(),
		})
	}
	//如果获取到的文件状态是文件夹
	if root.IsDir() {
		files, err := ioutil.ReadDir(rootpath)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": err.Error(),
			})
		}

		if path != "/" {
			pathPart := strings.Split(path, "/")
			if len(pathPart) > 0 {
				pathPart = pathPart[:len(pathPart)-1]
			}
			//定义上一级的父文件夹
			parentPath := strings.Join(pathPart, "/")
			parentFile := LocalFile{
				Name:     "..",
				Path:     "explorer?path=" + url.QueryEscape(parentPath), //对字符串转码，进行安全的url查询
				Size:     "",
				IsFolder: true,
				ModiTime: "",
			}
			localFiles = append(localFiles, parentFile)
			path = strings.Trim(path, "/") + "/" //把首尾的斜杠去掉，最后再加一个。因为path=那不能有斜杠
		} else {
			path = ""
		}
		//返回单个文件信息
		for _, fi := range files {
			file := LocalFile{
				Name:     fi.Name(),
				Path:     "explorer?path=" + url.QueryEscape(path+fi.Name()),
				Size:     Bytes2Size(fi.Size()),
				IsFolder: fi.Mode().IsDir(),
				ModiTime: fi.ModTime().String()[:19],
			}
			if file.IsFolder { //如果该文件是个文件夹
				localFiles = append(localFiles, file)
			} else {
				localf = append(localf, file)
			}
		}
		//合并所有文件夹和文件
		localFiles = append(localFiles, localf...)
		ctx.HTML(http.StatusOK, "explorer.html", gin.H{
			"message": "",
			"files":   localFiles,
		})
	} else {
		ctx.File(filepath.Join(LocalRoot, path))
	}
}
