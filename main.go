package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

type Item struct {
	Name        string `json:"name" bson:"name"`
	CreatedTime int64  `json:"created_time" bson:"created_time"` // 热搜数据保存到数据库的时间
	Type        int    `json:"type" bson:"type"`                 // 热搜数据的类型，目前只有两种： 0：热搜榜 1:要闻榜（标签榜）
	Index       int    `json:"index" bson:"index"`               // 数据的序号，热搜榜第一个数据是置顶数据，有51条。要闻榜50条
}

type Time struct {
	CreatedTime int64 `json:"created_time" bson:"created_time"` // 热搜数据保存到数据库的时间
}

func init() {
	initDB()
}

func main() {
	r := gin.Default()
	r.Use(Cors())

	r.GET("/api/list", func(c *gin.Context) {
		timeline := []Time{}
		err := timeCol.Find(ctx, bson.M{}).Sort("-created_time").All(&timeline)
		if err != nil {
			log.Fatal(err.Error())
		}

		created_time := c.DefaultQuery("created_time", strconv.FormatInt(timeline[0].CreatedTime, 10))

		created_time2, err := strconv.ParseInt(created_time, 10, 64)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("created_time2:", created_time2)

		items := getData(created_time2)

		// 转换时间线数据格式
		var timeArray []int64
		for _, item := range timeline {
			timeArray = append(timeArray, item.CreatedTime)
		}

		c.JSON(200, gin.H{
			"timeline": timeArray,
			"data":     items,
		})

	})
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "123123123")
	})
	r.Run()
}

// 时间戳
func getData(created_time int64) []Item {
	batch := []Item{}
	recordCol.Find(ctx, bson.M{"created_time": created_time}).All(&batch)
	return batch
}
