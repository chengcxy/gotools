package roboter

import (
	"github.com/chengcxy/gotools/configor"
)

type Roboter interface{
	SendMsg(content string) (string,error)
	GetPayload(content string)([]byte,error)

}



func GetRoboter(robotType string,config *configor.Config)Roboter{
	if robotType == "dingding"{
		return NewDingTalkRoboter(config)
	}else{
		return NewWechatRoboter(config)
	} 
}

