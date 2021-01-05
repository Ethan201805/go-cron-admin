package worker

import "go-cron/common"

type Scheduler struct {
	JobEventChan chan *common.JobEvent //etcd 任务事件队列
	JobPlanTable map[string] *common.JobSchedulePlan//任务调度计划表
}

//单例
var(
	G_scheduler *Scheduler
)

//调度协程
func (scheduler *Scheduler) scheduleLoop(){
	for{
		select {
		case JobEvent := <- scheduler.JobEventChan: //监听任务变化
			//对内存中任务表进行增删改查
			scheduler.HandleJobEvent(JobEvent)
		}
	}
}

//处理任务
func (scheduler *Scheduler) HandleJobEvent(jobEvent *common.JobEvent){
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE://保存任务事件
		//构建任务事件计划，并加入任务事件计划表
		plan, err := common.BuildJobSchedulePlan(jobEvent.Job)
		if err != nil{
			scheduler.JobPlanTable[jobEvent.Job.Name] = plan
		}
	case common.JOB_EVENT_DELETE://删除任务事件
		//任务计划表中存在则删除
		_,ok := scheduler.JobPlanTable[jobEvent.Job.Name]
		if ok{
			delete(scheduler.JobPlanTable,jobEvent.Job.Name)
		}
	}
}

//推送任务事件
func (Scheduler *Scheduler) PushScheduleEvent(event *common.JobEvent){
	Scheduler.JobEventChan <- event
}

//初始化调度器
func InitScheduler(){
	//赋值单例
	G_scheduler = &Scheduler{
		JobEventChan: make(chan *common.JobEvent,1000),
		JobPlanTable: make(map[string] *common.JobSchedulePlan),
	}
	//启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

