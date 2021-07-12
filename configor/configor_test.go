
package configor

import (
	"fmt"
	"testing"
)

var Usage = `
cd $project_path/edmgo/configor
go test

stdout result:

map[mysql:map[charset:utf8 db:test host:localhost password:admin123 port:3306 user:root]]
PASS
ok  	edmgo/configor	0.463s
`

func TestNewConfig(t *testing.T){
	ConfigPath := "/Users/chengxinyao/config/"
	Env := "dev"
	config := NewConfig(ConfigPath,Env)
	//fmt.Println(config.Conf)
	fmt.Println(config.Get("taskmeta.conn"))
	c,_ := config.Get("job_meta_conf")
	fmt.Println(c)
	c,_ = config.Get("roboter.roboter_type")
	fmt.Println(c)


}




