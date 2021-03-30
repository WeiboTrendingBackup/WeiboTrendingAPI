package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type Item struct {
	Name        string `json:"name"`
	CreatedTime int64  `json:"created_time"` // 热搜数据保存到数据库的时间
	Type        int    `json:"type"`         // 热搜数据的类型，目前只有两种： 0：热搜榜 1:要闻榜（标签榜）
	Index       int    `json:"index"`        // 数据的序号，热搜榜第一个数据是置顶数据，有51条。要闻榜50条
}

// Handler serverless-functions 函数暴露
func Handler(w http.ResponseWriter, r *http.Request) {
	initDB()
	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}

var recordCol *qmgo.Collection
var ctx context.Context
var dbClient *qmgo.Client

var MONGODB_URI string

// 因为此函数依赖log，因此不能放在init函数中执行，需要后于log.go的init函数执行
func initDB() {
	ctx = context.Background()
	var err error
	dbClient, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
	if err != nil {
		log.Fatalln("连接 Mongo 数据库报错", err.Error())
	}

	db := dbClient.Database("weibo")
	recordCol = db.Collection("record")

	// 各个数据的索引key，必须有一个非重复的数据，比如id。
	// Unique: 你的索引条件是否要求唯一。注意：是说整个【Key数组】匹配的结果是否唯一，而不是说单独的key是否唯一。
	// 每次修改索引之后，得手动删掉数据库的collection才能生效。（似乎只删除coll下面的 indexes 文件夹就行？
	recordCol.CreateOneIndex(ctx, options.IndexModel{Key: []string{"created_time", "type"}, Unique: false, Background: true})
}
