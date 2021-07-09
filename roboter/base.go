package roboter

import (
	"github.com/chengcxy/gotools/configor"
)

type Roboter interface{
	SendMsg(content string) (string,error)
	GetPayload(content string)([]byte,error)

}



func GetRoboter(robotType string,config *configor.Config)Roboter{
	if robotType == "weixin"{
		return NewWechatRoboter(config)
	}else{
		return NewDingTalkRoboter(config)
	} 
}


