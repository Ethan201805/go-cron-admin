package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct{
	ApiPort int	`json:"apiPort"`
	ApiReadTimeOut int `json:"apiReadTimeOut"`
	ApiWriteTimeOut int `json:"apiWriteTimeOut"`
	EtcdEndPoints []string `json:"etcdEndpoints"`
	EtcdDailTimeOut int `json:"etcdDailTimeOut"`
	StaticDir string `json:"staticDir"`
}

var(
	G_config *Config
)
//加载配置
func InitConfig(filename string) (err error){
	var(
		content []byte
		conf Config
	)
	//1 读取配置文件
	if content, err = ioutil.ReadFile(filename);err != nil{
		return
	}
	//2 反序列化
	if err = json.Unmarshal(content,&conf);err != nil{
		return
	}
	//3 赋值单例
	G_config = &conf

	return
}

