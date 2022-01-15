package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

//定义文件信息结构体 数据库字段设计
type File struct {
	Id              int    `json:"id"`
	Filename        string `json:"filename" gorm:"type:string"` //gorm框架指定数据库字段的数据类型
	Description     string `json:"description" gorm:"type:string"`
	Uploader        string `json:"uploader" gorm:"type:string unique"` //unique关键字表示唯一，路径不能重复
	Path            string `json:"path" gorm:"type:string"`
	Time            string `json:"time" gorm:"type:string"`
	DownloadCounter int    `json:"download_counter" gorm:"type:int"`
	IsLocalFile     bool   `json:"is_local_file" gorm:"type:bool"`
}

//实体化一个新的数据库
var DB *gorm.DB

//执行文件信息查询及关键字检索
func Query(query string) ([]*File, error) {
	var files []*File
	var err error
	//防止大小写敏感
	query = strings.ToLower(query)
	//使用gorm的where关键字进行查询定义，从files中进行查询，并将结果按id降序排列
	err = DB.Where("filename LIKE ? or description LIKE ? or uploader LIKE ? or time LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Order("id desc").Find(&files).Error
	return files, err
}
