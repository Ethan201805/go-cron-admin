package api

import (
	"github.com/gin-gonic/gin"
	"go-cron/common"
	"go-cron/master"
	"go-cron/tools"
	"go-cron/tools/app"
	"net/http"
)

//任务的http接口
type ApiServer struct {
	httpServer *http.Server
}
//配置单例
var (
	G_apiServer *ApiServer
)

//任务保存到ETCD POST {"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func SaveJob(c *gin.Context){
	var (
		job common.Job
		oldJob *common.Job
		err error
	)
	err = c.ShouldBindJSON(&job)
	//异常直接返回
	tools.HasError(err, "", 500)
	//保存到ETCD,调用JobMgr对象内SaveJob方法
	oldJob, err = master.G_jobMgr.SaveJob(&job)
	//异常直接返回
	tools.HasError(err, "", -1)
	//返回正常结果
	app.OK(c, oldJob, "success")
}


//从ETCD删除任务 POST name = "job1"
func DeleteJob(c *gin.Context){
	var(
		//oldJob *common.Job
		name string
		err error
	)
	//获取参数
	name = c.Param("name")
	//从etcd中删除指定name
	_,err = master.G_jobMgr.DelJob(name)
	//异常panic
	tools.HasError(err,"",-1)
	//返回正常结果
	app.OK(c,name,"success")
}

//任务列表
func JobList(c *gin.Context){
	//从etcd中获取数据列表
	list,err := master.G_jobMgr.JobList()
	//异常panic
	tools.HasError(err,"",-1)
	//返回正常结果
	app.OK(c,list,"success")
}

//kill任务
func KillJob(c *gin.Context){
	//获取参数
	name := c.Param("name")
	//从etcd中新增kill,通知worker去kill
	err := master.G_jobMgr.KillJob(name)
	//异常panic
	tools.HasError(err,"",-1)
	//返回正常结果
	app.OK(c,nil,"success")



}

//func HandleJobSave(resp http.ResponseWriter,req *http.Request){
//	var(
//		err error
//		postParam string
//		job common.Job
//		oldJob *common.Job
//		bytes []byte
//
//	)
//	//1 解析表单
//	if err = req.ParseForm();err != nil{
//		goto ERR
//	}
//	//2 获取参数
//	postParam = req.PostForm.Get("job")
//	//3 反序列化成job对象
//	if err = json.Unmarshal([]byte(postParam),&job);err != nil {
//		goto ERR
//	}
//	//4 保存到ETCD,调用JobMgr对象内SaveJob方法
//	if oldJob,err = master.G_jobMgr.SaveJob(&job);err != nil{
//		goto ERR
//	}
//	//5 返回正常应答
//	if bytes,err = common.BuildResponse(20000,"success",oldJob);err == nil{
//		resp.Write(bytes)
//	}
//
//	return
//ERR:
//	//6 返回异常应答
//	if bytes,err = common.BuildResponse(50000,err.Error(),nil);err == nil{
//		resp.Write(bytes)
//	}
//}

//func HandleJobDelete(resp http.ResponseWriter,req *http.Request) {
//	var(
//		err error
//		bytes []byte
//		name string
//		oldJob *common.Job
//	)
//	//1 解析表单
//	if err = req.ParseForm();err!=nil{
//		goto ERR
//	}
//	//2 获取参数
//	name = req.PostForm.Get("name")
//
//	//3 从etcd中删除指定name
//	if oldJob,err = master.G_jobMgr.DelJob(name);err!=nil{
//		goto ERR
//	}
//	//4 返回正常应答
//	if bytes,err = common.BuildResponse(20000,"success",oldJob);err == nil{
//		resp.Write(bytes)
//	}
//	return
//ERR:
//	//5 返回异常应答
//	if bytes,err = common.BuildResponse(50000,err.Error(),nil);err == nil{
//		resp.Write(bytes)
//	}
//}

//任务列表
//func HandleJobList(resp http.ResponseWriter,req *http.Request) {
//	// 从etcd中获取列表
//	list, err := master.G_jobMgr.JobList()
//	if err != nil{
//		goto ERR
//	}
//	//返回正常应答
//	if bytes,err := common.BuildResponse(20000,"success",list);err == nil{
//		resp.Write(bytes)
//	}
//	return
//ERR:
//	//返回异常应答
//	if bytes,err := common.BuildResponse(50000,err.Error(),nil);err == nil{
//		resp.Write(bytes)
//	}
//}
//杀死任务
//func HandleJobKill(resp http.ResponseWriter,req *http.Request) {
//	var(
//		err error
//		name string
//		bytes []byte
//	)
//	// 解析表单
//	if err = req.ParseForm();err!=nil{
//		goto ERR
//	}
//	// 获取参数
//	name = req.PostForm.Get("name")
//
//	// etcd中新增kill 操作
//	if err = master.G_jobMgr.KillJob(name);err!=nil{
//		goto ERR
//	}
//	//返回正常应答
//	if bytes,err = common.BuildResponse(20000,"success",nil);err == nil{
//		resp.Write(bytes)
//	}
//	return
//ERR:
//	//返回异常应答
//	if bytes,err = common.BuildResponse(50000,err.Error(),nil);err == nil{
//		resp.Write(bytes)
//	}
//}

//初始化服务
//func InitApiServer()(err error){
//	var (
//		mux *http.ServeMux
//		listener net.Listener
//		httpServer *http.Server
//		staticDir http.Dir
//		staticHandle http.Handler
//	)
//	//配置路由
//	mux = http.NewServeMux()
//	mux.HandleFunc("/job/save", handleJobSave)
//	mux.HandleFunc("/job/delete", handleJobDelete)
//	mux.HandleFunc("/job/list", handleJobList)
//	mux.HandleFunc("/job/kill", handleJobKill)
//
//	//配置静态文件目录
//	staticDir = http.Dir(master.G_config.StaticDir)
//	staticHandle = http.FileServer(staticDir)
//	mux.Handle("/",http.StripPrefix("/",staticHandle))
//
//	//启动tcp监听
//	listener,err = net.Listen("tcp",":"+strconv.Itoa(master.G_config.ApiPort));if err != nil {
//		return
//	}
//	//创建一个http服务
//	httpServer = &http.Server{
//		ReadTimeout: time.Duration(master.G_config.ApiReadTimeOut)*time.Millisecond,
//		WriteTimeout: time.Duration(master.G_config.ApiWriteTimeOut)*time.Millisecond,
//		Handler: mux,
//	}
//	//单例赋值
//	G_apiServer = &ApiServer{
//		httpServer: httpServer,
//	}
//	//启动服务
//	go httpServer.Serve(listener)
//
//	return
//}


