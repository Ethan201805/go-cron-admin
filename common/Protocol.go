package common

import (
	"encoding/json"
	"strings"
	"time"
	"github.com/gorhill/cronexpr"
)

//定时任务
type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//http接口应答
type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type JobEvent struct {
	EventType int
	Job *Job
}

type JobSchedulePlan struct {
	Job *Job
	Expr *cronexpr.Expression
	NextTime time.Time
}

//应答方法
func BuildResponse(code int,message string,data interface{})(resp []byte,err error){
	//1 定义一个Response
	var(
		response Response
	)
	//赋值
	response.Code = code
	response.Message = message
	response.Data = data
	//序列化
	resp,err = json.Marshal(response)
	return
}

//反序列化
func UnpackJob(value []byte)(j *Job,err error){
	var(
		job = &Job{}
	)
	if err = json.Unmarshal(value,job);err!=nil{
		return
	}
	j = job
	return
}

//提取任务名
func TrimJobName(key string)(name string){
	return strings.TrimPrefix(key,JOB_SAVE_DIR)
}

//任务变化事件：1 更新 2删除
func BuildJobEvent(eventType int,job *Job)(eventJob *JobEvent){
	return &JobEvent{
		EventType: eventType,
		Job: job,
	}
}

//构建任务事件计划
func BuildJobSchedulePlan(job *Job)(jobSchedulePlan *JobSchedulePlan,err error){
	//解析表达式
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil{
		return
	}
	//生成任务调度计划对象
	jobSchedulePlan =  &JobSchedulePlan{
		Job:job,
		Expr: expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}