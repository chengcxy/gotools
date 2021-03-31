package roboter

import (
	"fmt"
	"flag"
	"encoding/json"
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
	"github.com/chengcxy/gotools/configor"

)



type Roboter struct {
	Config *configor.Config
}


func NewRoboter(config *configor.Config)(*Roboter){
	r := &Roboter{
		Config:config,
	}
	return r
}


func (r *Roboter) GetPayload(content string) (string,map[string]interface{}){
	data := make(map[string]interface{})
	text := make(map[string]string)
	v,ok := r.Config.Get("roboter")
	if !ok {
		panic ("roboter is not exists")
	}
	rc := v.(map[string]interface{})
	hook_keyword := rc["hook_keyword"].(string)
	text["content"] = hook_keyword + content
	at := make(map[string]interface{})
	token := rc["token"].(string)
	at["atMobiles"] = rc["atMobiles"].([]interface{})
	at["isAtAll"] = rc["isAtAll"].(bool)
	data["msgtype"] = "text"
	data["text"] = text
	data["at"] = at
	return token,data
}
func (r *Roboter) SendMsg(content string)string{
	contentType := "application/json;charset=utf-8"
	token,data := r.GetPayload(content)
	api := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s",token)
	fmt.Println(api)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
    jsonStr, _ := json.Marshal(data)
    resp, err := client.Post(api, contentType, bytes.NewBuffer(jsonStr))
    if err != nil {
        panic("send dingtalk msg err")
    }
    defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
    return string(result)

}


var Usage = `
go run main.go  -config_path /data/config -env test 
`


//命令行参数
var ConfigPath = flag.String("config_path","/data/", "配置文件父目录")
var Env = flag.String("env","local", "json配置文件名,不带扩展名")
var Msg = flag.String("msg","test roboter", "消息")
//被其他模块调用的Start方法
func InitRoboter(){
	flag.Parse()
	config_path := *ConfigPath
	env := *Env
	msg := *Msg
	//解析json配置文件
	config := configor.NewConfig(config_path,env)
	r := NewRoboter(config)
	r.SendMsg(msg)
}
