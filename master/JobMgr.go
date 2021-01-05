package master

import (
	"context"
	"encoding/json"
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
}
//配置单例
var (
	G_jobMgr *JobMgr
)
//初始化
func InitJobMgr()(err error){
	var(
		conf clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
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

	//赋值单例
	G_jobMgr = &JobMgr{
		client:client,
		kv:kv,
		lease:lease,
	}
	return
}
//保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job)(oldJob *common.Job,err error){
	//保存到/cron/jobs/任务名
	var(
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldValue common.Job
	)
	//etcd 保存key
	jobKey = common.JOB_SAVE_DIR+job.Name
	//etcd 保存value,对象序列化成json
	if jobValue,err = json.Marshal(job);err != nil{
		return
	}
	//保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV());err != nil{
		fmt.Println(err.Error())
		return
	}
	//如果是更新，返回旧值
	if putResp.PrevKv != nil{
		//对旧值做反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value,&oldValue);err != nil{
			fmt.Println(err.Error())
			err = nil
			return
		}
		oldJob = &oldValue
	}
	return
}
//删除任务
func (jobMgr *JobMgr) DelJob(name string)(oldJob *common.Job,err error){
	var(
		jobKey string
		resp *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	//etcd中保存任务的key
	jobKey = common.JOB_SAVE_DIR + name
	//从etcd中删除
	if resp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil{
		return
	}
	//返回删除的值
	if len(resp.PrevKvs) != 0{
		//反序列化旧值
		if err = json.Unmarshal(resp.PrevKvs[0].Value,&oldJobObj);err != nil{
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

//获取任务列表
func (jobMgr *JobMgr) JobList()(jobList []*common.Job ,err error){
	//定义列表目录
	var(
		jobKey string
		getResp *clientv3.GetResponse
		kvpair *mvccpb.KeyValue
	)
	jobKey = common.JOB_SAVE_DIR
	//获取任务列表
	if getResp, err = jobMgr.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err!=nil{
		return
	}
	//初始化返回数组
	jobList = make([]*common.Job,0)

	//遍历内部k-v值并反序列化成结构体
	for _,kvpair = range getResp.Kvs{
		var job = common.Job{}
		if err = json.Unmarshal(kvpair.Value, &job);err != nil{
			err = nil
			continue
		}
		jobList = append(jobList,&job)
	}

	return
}

//etcd 通知worker 杀死任务
func (jobMgr *JobMgr) KillJob(name string)(err error){
	//生成killKey
	killKey := common.JOB_KILL_DIR+name
	//生成1s租约
	grant, err := jobMgr.lease.Grant(context.TODO(), 1)
	if err != nil{
		return
	}
	//带着租约往/cron/kill目录插入需要kill的任务,让worker监听到一次put操作
	_, err = jobMgr.kv.Put(context.TODO(), killKey, "", clientv3.WithLease(grant.ID))
	if err != nil{
		return
	}
	return
}
