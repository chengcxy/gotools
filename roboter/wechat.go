package roboter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/chengcxy/gotools/configor"
)

/*
微信报警 json配置文件 
"roboter": {
	"token": "token",
	"isAtAll": false,
	
}
*/




//钉钉机器人post请求接口地址
var WechatBaseApi = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

type WechatRoboter struct{
	Token string
}


func (dt WechatRoboter)SendMsg(content string)(string,error){
	payload,err := dt.GetPayload(content)
	if err != nil{
		log.Println("get wechat payload message error,",err)
		return fmt.Sprintf("content:%s transfer bytes error",content),err
	}
	contentType := "application/json"
	api := fmt.Sprintf(WechatBaseApi,dt.Token)
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
    resp, err := client.Post(api, contentType, bytes.NewBuffer(payload))
    if err != nil {
        panic("send dingtalk msg err")
    }
    defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
    return string(result),err
}

func (dt WechatRoboter)GetPayload(content string)([]byte,error){
	data := make(map[string]interface{})
	data["msgtype"] = "text"
	text := make(map[string]string)
	text["content"] = fmt.Sprintf("%s",content)
	data["text"] = text
	return json.Marshal(data)

}


func NewWechatRoboter(config *configor.Config) Roboter{
	v,ok := config.Get("roboter")
	if !ok {
		panic("roboter key is not in json config_file")
	}
	rc := v.(map[string]interface{})
	token := rc["token"].(string)
	return WechatRoboter{
		Token:token,
	}
}


