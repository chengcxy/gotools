
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
	ConfigPath := "/Users/chengxinyao/go/src/edmgo/JsonConfigFiles"
	Env := "local"
	Config := NewConfig(ConfigPath,Env)
	fmt.Println(Config.Conf)

}
