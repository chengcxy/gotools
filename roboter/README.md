
# 报警机器人如何使用 

- 配置文件 比如dev.json 添加roboter的key 钉钉机器人配置规则

```json
{
    "roboter": {
        "token": "token",
        "atMobiles": [
            "$mobile"
        ],
        "isAtAll": false,
        "hook_keyword": "任务报警",
        "roboter_type": "dingding"
     }
}
```


```go
package main


import(
    "github.com/chengcxy/gotools/configor"
    "github.com/chengcxy/gotools/roboter"

)


func main(){
    ConfigPath := "/Users/chengxinyao/config"
	Env := "dev"
	config := configor.NewConfig(ConfigPath,Env)
	robot := roboter.GetRoboter("dingding",config)
	robot.SendMsg("get dingding robot")

}

```

