package roboter

import (
	"github.com/chengcxy/gotools/configor"
)

type Roboter interface{
	SendMsg(content,mobile string) (string,error)
	GetPayload(content,mobile string)([]byte,error)

}




func GetRoboter(config *configor.Config)Roboter{
	c,_ := config.Get("roboter")
	robotType := c.(map[string]interface{})["roboter_type"].(string)
	if robotType == "weixin"{
		return NewWechatRoboter(config)
	}else{
		return NewDingTalkRoboter(config)
	} 
}


