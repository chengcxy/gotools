
package configor

import (
	"path"
	"github.com/chengcxy/gotools/utils"
	"strings"
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

func getConfig(conf interface{},key string)(interface{},bool){
	mp := conf.(map[string]interface{})
	value,ok := mp[key]
	if ok {
		return value,ok
	}
	isReparse := strings.Contains(key,".")
	if !isReparse{
		value,ok := mp[key]
		if ok {
			return value,ok
		}
		return nil,false	
	}
	keys := strings.Split(key,".") 
	for len(keys) > 0 {
		key = keys[0]
		keys = keys[1:]
		strKeys := strings.Join(keys,".")
		value,ok := mp[key]
		if ok{
			return getConfig(value,strKeys)
		}
		return nil,false
	}
	return nil,false
}

// json key like "a.b.c" 
func(c *Config) Get(key string)(interface{},bool){
	return getConfig(c.Conf,key)
}

