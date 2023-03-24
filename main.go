package main

import (
	"douyin/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	go controller.RunMessageServer()
	controller.InitDb()
	r := gin.Default()

	initRouter(r)
	controller.InitRedis()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
