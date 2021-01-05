package router

import (
	"github.com/gin-gonic/gin"
	"go-cron/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/",func(c *gin.Context) {
		c.String(200, "Api Running")
	})

	//注册中间件
	middleware.InitMiddleware(r)
	// 注册业务路由
	InitJobRouter(r)

	return r
}

