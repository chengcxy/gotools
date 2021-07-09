
package roboter

import (
	"fmt"
	"testing"
	"github.com/chengcxy/gotools/configor"
)


/* 测试通过
{"errcode":0,"errmsg":"ok"} <nil>
{"errcode":93000,"errmsg":"invalid webhook url, hint: [1625817226_237_59554d8be5fdec74ad83e88e99b6c607], from ip: 58.33.178.162, more info at https://open.work.weixin.qq.com/devtool/query?e=93000"} <nil>
PASS
ok  	github.com/chengcxy/gosqlboy/roboter	1.260s
*/


//钉钉机器人
func TestNewDingTalkRoboter(t *testing.T){
	ConfigPath := "/Users/chengxinyao/config"
	Env := "dev"
	config := configor.NewConfig(ConfigPath,Env)
	robot := NewDingTalkRoboter(config)
	fmt.Println(robot.SendMsg("测试"))
}

//微信机器人
func TestNewWechatRoboter(t *testing.T){
	ConfigPath := "/Users/chengxinyao/config"
	Env := "dev"
	config := configor.NewConfig(ConfigPath,Env)
	robot := NewWechatRoboter(config)
	fmt.Println(robot.SendMsg("测试"))
}

//微信机器人
func TestGetRoboter(t *testing.T){
	ConfigPath := "/Users/chengxinyao/config"
	Env := "dev"
	config := configor.NewConfig(ConfigPath,Env)
	robot := GetRoboter("dingding",config)
	robot.SendMsg("get dingding robot")
}