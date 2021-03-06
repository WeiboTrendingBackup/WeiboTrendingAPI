package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qiniu/qmgo"
)

var recordCol *qmgo.Collection
var timeCol *qmgo.Collection
var ctx context.Context
var dbClient *qmgo.Client

var MONGODB_URI string
var ENV string

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("加载 .env 文件失败", err.Error())
	}

	MONGODB_URI = os.Getenv("MONGODB_URI")
	ENV = os.Getenv("ENV")
	// 开发环境数据库链接
	if ENV == "dev" {
		MONGODB_URI = "mongodb://localhost:27017"
	}
}

// 因为此函数依赖log，因此不能放在init函数中执行，需要后于log.go的init函数执行
func initDB() {
	loadConfig()

	ctx = context.Background()
	var err error
	dbClient, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: MONGODB_URI})
	if err != nil {
		log.Fatalln("连接 Mongo 数据库报错", err.Error())
	}

	db := dbClient.Database("weibo")
	recordCol = db.Collection("record")
	timeCol = db.Collection("time")
}
