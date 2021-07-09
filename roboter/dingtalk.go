
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
钉钉报警json配置文件  如果机器人设置了信息关键字 hook_keyword填写该关键字 否则无法推送成功
比如 我的机器人设置了"任务"2个字 只有消息里面含有任务字样才可以发送
"roboter": {
	"token": "token",
	"atMobiles": [
		"$mobile"
	],
	"isAtAll": false,
	"hook_keyword": "任务报警"
}

*/

//钉钉机器人post请求接口地址
var DingTalkBaseApi = "https://oapi.dingtalk.com/robot/send?access_token=%s"

type DingTalkRoboter struct{
	Token string
	AtMobiles []interface{}
	Hookeyword string
	IsAtall bool
}


func (dt DingTalkRoboter)SendMsg(content string)(string,error){
	payload,err := dt.GetPayload(content)
	if err != nil{
		log.Println("get dingtalk payload message error,",err)
		return fmt.Sprintf("content:%s transfer bytes error",content),err
	}
	contentType := "application/json;charset=utf-8"
	api := fmt.Sprintf(DingTalkBaseApi,dt.Token)
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

func (dt DingTalkRoboter)GetPayload(content string)([]byte,error){
	data := make(map[string]interface{})
	data["msgtype"] = "text"
	at := make(map[string]interface{})
	at["atMobiles"] = dt.AtMobiles
	at["isAtAll"] = dt.IsAtall
	data["at"] = at
	text := make(map[string]string)
	text["content"] = fmt.Sprintf("%s:\n%s",dt.Hookeyword,content)
	data["text"] = text
	return json.Marshal(data)

}



func NewDingTalkRoboter(config *configor.Config) Roboter{
	v,ok := config.Get("roboter")
	if !ok {
		panic("roboter key is not in json config_file")
	}
	rc := v.(map[string]interface{})
	token := rc["token"].(string)
	atMobiles := rc["atMobiles"].([]interface{})
	hook_keyword := rc["hook_keyword"].(string)
	isAtAll := rc["isAtAll"].(bool)
	return DingTalkRoboter{
		Token:token,
		AtMobiles:atMobiles,
		Hookeyword:hook_keyword,
		IsAtall:isAtAll,
	}
}


