package main

import (
	"douyin/controller"
	"douyin/dao"

	"github.com/gin-gonic/gin"
)

func main() {
	go controller.RunMessageServer()
	dao.InitDb()
	controller.InitDb()
	r := gin.Default()

	initRouter(r)
	controller.InitRedis()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
