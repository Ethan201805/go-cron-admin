package worker

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go-cron/common"
	"time"
)
//任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	watcher clientv3.Watcher
}
//配置单例
var (
	G_jobMgr *JobMgr
)
//监听任务变化
func (JobMgr *JobMgr) watchJobs()(err error){
	var(
		jobEvent *common.JobEvent
	)
	//1 get一下/cron/jobs下所有任务，并货值当前集群的revision
	getResp, err := JobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())
	if err != nil{
		return
	}
	for _,kvpair := range getResp.Kvs{
		//反序列化
		job, err := common.UnpackJob(kvpair.Value)
		if err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//todo 把这个job同步给scheduler调度携程
		}
	}
	//2 从该revision向后监听变化事件
	 go func() {
	 	//从get时刻的后续版本开始监听
	 	watchStartRevision := getResp.Header.Revision+1
	 	//监听/cron/jobs/目录的后续变化
		 watchChan := G_jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		 //处理监听事件
		 for watchResp := range watchChan{
		 	for _,watchEvent := range watchResp.Events{
				switch watchEvent.Type{
				case mvccpb.PUT:
					//todo 反序列化job
					job, err := common.UnpackJob(watchEvent.Kv.Value)
					if err == nil{
						continue
					}
					//构建一个事件
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					//获取任务名
					jobName := common.TrimJobName(string(watchEvent.Kv.Key))
					//获取一个job
					job := &common.Job{Name:jobName}
					//构建一个事件
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				//todo 推一个事件给scheduler
				G_scheduler.PushScheduleEvent(jobEvent)
			}
		 }
	 }()
	return
}

//初始化
func InitJobMgr()(err error){
	var(
		conf clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		watcher clientv3.Watcher
	)
	//初始化配置
	conf = clientv3.Config{
		Endpoints:G_config.EtcdEndPoints,//集群地址
		DialTimeout:time.Duration(G_config.EtcdDailTimeOut)*time.Millisecond,
	}
	//建立连接
	if client,err = clientv3.New(conf);err !=nil {
		fmt.Println(err)
		return
	}
	//得到kv和lease的api子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client:client,
		kv:kv,
		lease:lease,
		watcher:watcher,
	}
	//启动任务监听
	G_jobMgr.watchJobs()

	return
}
