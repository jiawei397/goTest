package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
)

type File struct {
	_id      string `bson:"id"`
	count    int
	deleted  bool
	name     string
	ossName  string
	url      string
	size     int
	mimeType string
	userId   string
	// createTime date
	// modifyTime date
}

func main() {
	// 设置客户端连接配置
	session, err := mgo.Dial("mongodb://192.168.21.176:27018")

	// 连接到MongoDB
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// 指定获取要操作的数据集
	c := session.DB("qdownload").C("files")

	fmt.Println(c.Count())

	// var files []map[string]interface{}
	var files []File
	c.Find(map[string]string{"name": "node-v14.3.0.pkg"}).All(&files)
	fmt.Println(files)

	// file := File{}
	var file map[string]interface{}
	c.Find(bson.M{"name": "node-v14.3.0.pkg"}).One(&file)
	fmt.Println(file["url"])
}
