
package configor

import (
	"path"
	"github.com/chengcxy/gotools/utils"
)

type Config struct{
	ConfigPath string
	Env string
	Conf map[string]interface{}
}



func NewConfig(ConfigPath,Env string)(*Config){
	c := &Config{
		ConfigPath:ConfigPath,
		Env:Env,
	}
	conf,_ := c.getJsonConfig()
	c.Conf = conf
	return c
}


func (c *Config)getJsonConfig()(map[string]interface{},error){
	json_file := path.Join(c.ConfigPath,c.Env + ".json")
	m,err := utils.ParseJsonFile(json_file)
	if err != nil{
		panic(err)
	}
	return m,nil
}

func(c *Config) Get(key string)(interface{},bool){
	value,ok := c.Conf[key]
	if ok {
		return value,ok
	}
	return nil,false

}

