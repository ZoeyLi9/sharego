package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"strings"
)

//定义文件信息结构体 数据库字段设计
type File struct {
	Id              int    `json:"id"`
	Filename        string `json:"filename" gorm:"type:string"` //gorm框架指定数据库字段的数据类型
	Description     string `json:"description" gorm:"type:string"`
	Uploader        string `json:"uploader" gorm:"type:string"`
	Link            string `json:"link" gorm:"type:string unique"` //unique关键字表示唯一，路径不能重复
	Time            string `json:"time" gorm:"type:string"`
	DownloadCounter int    `json:"download_counter" gorm:"type:int"`
	IsLocalFile     bool   `json:"is_local_file" gorm:"type:bool"`
}

//用来浏览本地文件信息结构
type LocalFile struct {
	Name     string
	Path     string
	Size     string
	IsFolder bool
	ModiTime string
}

//实体化一个新的数据库
var DB *gorm.DB

//初始化数据库
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./.sharego.db") //打开一个数据库连接，若没有，则新建一个数据库文件
	if err == nil {                                  //如果没有报错
		DB = db
		db.AutoMigrate(&File{}) //增量转移
		return DB, err
	} else {
		log.Fatal(err)
	}
	return nil, err
}

//数据库插入记录的方法
func (file *File) Insert() error {
	var err error
	err = DB.Create(file).Error //将信息插入数据库中
	return err
}

//数据库删除记录的方法
func (file *File) Delete() error {
	var err error
	err = DB.Delete(file).Error
	//删除本地文件
	if !file.IsLocalFile {
		err = os.Remove("." + file.Link)
	}
	return err
}

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