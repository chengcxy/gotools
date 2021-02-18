package utils


import (
	"errors"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"path"
)



func ParseJsonFile(json_file string)(map[string]interface{},error){
	file,err := os.Open(json_file)
	defer file.Close()
	if err != nil{
		return nil,err
	}
	data,_ := ioutil.ReadAll(file)
	var f interface{}
	json.Unmarshal(data, &f)
	m := f.(map[string]interface{})
	return m,nil

}

//检验代码里是否包含关键字
type CheckCodeContainsWdResult struct{
	File string
	Wd   string
	Flag bool
}


//递归遍历目录 寻找以.go/.py/.java ... 为扩展名结尾的文件
func RecursionPath(folder_path,end_flag string)(vs []string){
	f,_ := os.Stat(folder_path)
	is_dir := f.IsDir() 
	exist_flag := strings.HasSuffix(folder_path, end_flag)
	if !is_dir && exist_flag{
		vs = append(vs,folder_path)
	}else {
		if is_dir{
			fileInfo, _ := ioutil.ReadDir(folder_path)
			for _, v := range fileInfo {
				_v := RecursionPath(path.Join(folder_path,v.Name()),end_flag)
				vs = append(vs,_v ...)
			}
		}
	}
  return vs
}

//检查单文件是否存在关键字
func CheckWordExists(file,wd string,ch chan *CheckCodeContainsWdResult)(_flag bool){
	data, _ := ioutil.ReadFile(file)
	s := string(data)
	_flag = strings.Contains(s,wd)
	r := &CheckCodeContainsWdResult{
		File:file,
		Wd:wd,
		Flag:_flag,
	}
	ch <- r
	return _flag
}


//对表数据 根据batch的数据进行切分 切片套切片 [[0,10000],[10000,20000]...]
func SplitBatch(start,end,batch int)([][]int,int){
	l := make([][]int,0)
	process_num := 0
	for start < end{
		p := make([]int,2)
		p[0] = start
		_end := start + batch
		if _end >= end{
			_end = end
		}
		p[1] = _end
		l = append(l,p)
		process_num ++
		start = _end
	}
	return l,process_num
}

//根据表的总量 进行一个切分
func GetBatch(total_count int)(int){
	batch := 1000
	if total_count <= 50000{
		batch = 3000
	}else if total_count <= 100000 {
		batch = 10000
	}else if total_count <=200000{
		batch = 15000
	}else if total_count <=500000{
		batch = 25000
	}else{
		batch = 30000
	}
	return batch
}

//拼接mysql insert into sql语句
//如果 is_on_duplicate 带上 ON DUPLICATE KEY UPDATE
// fields=[a,b,c] uniqueindexs=[c] duplicate_keys=[a,b]

func GenInsertSql(db,table string,fields []string,uniqueindexs[]string,value_list_num int) (string,error){
	//拼接占位符 
	if value_list_num <0{
		return "",errors.New("value_list_num <0 is not allowed")
	}
	uniqueindexs_num := 0
	if uniqueindexs != nil{
		uniqueindexs_num = len(uniqueindexs)
	}

	duplicate_keys := make([]string,0)

	questions := make([]string,len(fields))
	
	for i:=0;i<len(fields);i++{
		questions[i] = "?"
		if uniqueindexs_num !=0 {
			for j:=0;j<uniqueindexs_num;j++{
				if fields[i] != uniqueindexs[j]{
					duplicate_keys = append(duplicate_keys,fmt.Sprintf("%s=values(%s)",fields[i],fields[i]))
				}
			}	
		}else{
			duplicate_keys = append(duplicate_keys,fmt.Sprintf("%s=values(%s)",fields[i],fields[i]))
		}
	}
	questions_str := "(" + strings.Join(questions,",")+")"
	
	all_questions := make([]string,value_list_num/len(fields))
	field_strs := strings.Join(fields,",")
	duplicate_keys_str := strings.Join(duplicate_keys,",")
	for i:=0;i<value_list_num/len(fields);i++{
		all_questions[i] = questions_str
	}
	all_questions_str := strings.Join(all_questions,",")
	sql := fmt.Sprintf("insert into %s.%s(%s)values %s ON DUPLICATE KEY UPDATE %s",db,table,field_strs,all_questions_str,duplicate_keys_str)
	return sql,nil
}