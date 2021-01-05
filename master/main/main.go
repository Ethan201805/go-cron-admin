package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"go-cron/master"
	"go-cron/router"
	"go-cron/tools"
	"runtime"
)

var(
	filePath string
)

func initEnv(){
	//设置线程数为机器核心数目
	runtime.GOMAXPROCS(runtime.NumCPU())
}
//初始化参
func initArgs(){
	//master -config ./worker.json
	flag.StringVar(&filePath,"config","./master/master.json","load worker.json")
	flag.Parse()
}

func main() {
	var(
		err error
		r *gin.Engine
	)
	//初始化路由
	r = router.InitRouter()
	//初始化线程
	initEnv()
	//初始化参数
	initArgs()
	//加载配置文件
	err = master.InitConfig(filePath)
	tools.HasError(err,"",-1)
	//任务管理器
	err = master.InitJobMgr()
	tools.HasError(err,"",-1)

	//启动服务
	r.Run(":8090")

}
