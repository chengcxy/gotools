package backends

import (
	"fmt"
	//"log"
	"time"
	"strconv"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


//查询表的列字段
var QUERY_TABLE_COLUMNS = `
select column_name,column_type,column_comment,column_key
from information_schema.columns
where table_schema="%s" and table_name="%s"
`
//查询表的唯一索引
var QUERY_UNIQ_INDEXS =  "show index from %s.%s where non_unique = 0" 


//mysql配置信息
type MysqlConfig struct{
	Host string
	User string
	Password string
	Port int
	Db string
	Charset string
	ConnUri string

}

func NewMysqlConfig(f interface{})(*MysqlConfig){
	m := f.(map[string]interface{})
	mc := &MysqlConfig{
		Host:m["host"].(string),
		User:m["user"].(string),
		Password:m["password"].(string),
		Db:m["db"].(string),
		Port:int(m["port"].(float64)),
		Charset:m["charset"].(string),
	}
	mc.ConnUri = mc.GetConnUri()
	return mc
}

//拼接url
func (mc *MysqlConfig) GetConnUri()(url string){
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",mc.User,mc.Password,mc.Host,mc.Port,mc.Db,mc.Charset)
	return url
}



//表的元数据信息 数据库名 表名 主键 最大值 最小值 切分的任务列表 所属的客户端
type TableMeta struct{
	OwnApp string `json:"own_app"`
	DbName string `json:"db_name"`
	TableName string `json:"table_name"`
	Pk string `json:"pk"`
	Fields []string `json:"fields"`
	UniqueIndexs []string `json:"unique_indexs"`
	MinId int `json:"min_id"`
	MaxId int `json:"max_id"`
	Batch int `json:"batch"` 
	HasPrimaryKey bool`json:"has_primary_key"` 
	TotalCount int `json:"total_count"`
}


//mysql客户端结构体
type MysqlClient struct{
	Config *MysqlConfig
	Db *sql.DB
}

//查询元数据
func (m *MysqlClient) GetTableMeta(own_app,db_name,table_name string)(*TableMeta){
	sql := fmt.Sprintf(QUERY_TABLE_COLUMNS,db_name,table_name)
	rows_list,_,err := m.Query(sql)
	if err != nil{
		panic("获取元数据失败")
	}

	fields := make([]string,len(rows_list))
	pk := ""
	has_pk := false
	for index,item := range rows_list{
		fields[index] = strings.ToLower(item["column_name"])
		if item["column_key"] == "PRI"{
			pk = strings.ToLower(item["column_name"])
			has_pk = true
		}
	}
	if !has_pk{
		panic("no pk")
	}
	min_id := m.GetMinId(db_name,table_name,pk)
	max_id := m.GetMaxId(db_name,table_name,pk)
	unique_indexs := m.GetUniqueIndexs(db_name,table_name)
	tm := &TableMeta{
		OwnApp:own_app,
		DbName:db_name,
		TableName:table_name,
		Pk:pk,
		Fields:fields,
		UniqueIndexs:unique_indexs,
		MinId:min_id,
		MaxId:max_id,
		HasPrimaryKey:has_pk,
	}
	return tm
}


func (m *MysqlClient) GetUniqueIndexs(db_name,table_name string) []string{
	sql := fmt.Sprintf(QUERY_UNIQ_INDEXS,db_name,table_name)
	rows_list,_,err := m.Query(sql)
	if err != nil{
		panic("获取唯一索引失败")
	}
	unique_indexs := make([]string,len(rows_list))
	if len(rows_list) ==0 {
		return unique_indexs
	}
	for index,item := range rows_list{
		unique_indexs[index] = strings.ToLower(item["column_name"])
	}
	
    return unique_indexs
}


//关闭数据库连接池
func (m *MysqlClient) Close() {
	m.Db.Close()
}


//根据主键id对表数据进行切分读取
func (m *MysqlClient) GetProcessSql(db_name,table_name,pk string) string{
	query := fmt.Sprintf("select * from %s.%s where %s>? and %s<=?",db_name,table_name,pk,pk)
	return query
}

//获取表最小值 
func (m *MysqlClient) GetTotalCount(db_name,table_name string) int{
	query := fmt.Sprintf("select count(1) as total from %s.%s ",db_name,table_name)
	rows,_,err := m.Query(query)
	if err != nil{
		panic("获取总数据量失败")
	}
	data := rows[0]["total"]
	total_count,_ := strconv.Atoi(data)
	return total_count
}



//获取表最小值 
func (m *MysqlClient) GetMinId(db_name,table_name,pk string) int{
	query := fmt.Sprintf("select %s from %s.%s order by %s limit 1",pk,db_name,table_name,pk)
	rows,_,err := m.Query(query)
	if err != nil{
		panic("获取最小id失败")
	}
	if len(rows) == 0{
		return 0
	}
	
	min_id_str := rows[0][pk]
	min_id,_ := strconv.Atoi(min_id_str)
	return min_id - 1
}
//获取表最大值 
func (m *MysqlClient) GetMaxId(db_name,table_name,pk string) int{
	query := fmt.Sprintf("select %s from %s.%s order by %s desc limit 1",pk,db_name,table_name,pk)
	rows,_,err := m.Query(query)
	if err != nil{
		panic("获取最大id失败")
	}
	if len(rows) == 0{
		return 0
	}
	max_id_str := rows[0][pk]
	max_id,_ := strconv.Atoi(max_id_str)
	return max_id
}









//mysql 客户端
func NewMysqlClient(mc *MysqlConfig)(*MysqlClient){
	Db,_ :=  sql.Open("mysql",mc.ConnUri)
	Db.SetConnMaxLifetime(time.Minute * 100)
	Db.SetMaxOpenConns(30)
	Db.SetMaxIdleConns(30)
	client := &MysqlClient{
		Config:mc,
		Db:Db,
	}
	return client
}

func (m *MysqlClient)Read(query string)([]map[string]string,[]string,error){
	return m.Query(query)
}


func (m *MysqlClient) Write(meta map[string]string,columns []string,is_create_table bool,datas []map[string]string)(int64,bool,int){
	status := 0
	to_db := meta["to_db"]
	to_table := meta["to_table"]
	write_batch,_ := strconv.Atoi(meta["write_batch"])
	write_mode := "insert"
	var num int64
	num = 0
	if len(columns)>0 && !is_create_table{
		stmts := make([]string,len(columns))
		schema := ""
		id_exists := false
		z_create_time_exists := false
		z_update_time_exists := false
		for _,col := range columns{
			if col == "id"{
				id_exists = true
			}
			if col == "z_create_time"{
				z_create_time_exists = true
			}
			if col == "z_update_time"{
				z_update_time_exists = true
			}
		}
		if id_exists{
			schema = fmt.Sprintf("create table if not exists %s.%s(add_id int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',",to_db,to_table)
		}else{
			schema = fmt.Sprintf("create table if not exists %s.%s(id int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',",to_db,to_table)
		}
		for index,col := range columns{
			stmts[index] = fmt.Sprintf(" %s varchar(255)",col)
		}
		schema += strings.Join(stmts,",")

		temp := make([]string,0)
		if !z_create_time_exists{
			temp = append(temp,"z_create_time datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'")
		}else{
			temp = append(temp,"add_z_create_time datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'")
		}
		if !z_update_time_exists{
			temp = append(temp,"z_update_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'")
		}else{
			temp = append(temp,"add_z_update_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'")
		}
		if !id_exists{
			temp = append(temp,"PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		}else{
			temp = append(temp,"PRIMARY KEY (add_id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		}
		schema += fmt.Sprintf(",%s",strings.Join(temp,","))
		m.Execute(schema)
		is_create_table = true
	}
	if len(datas) == 0{
		return 0,is_create_table,status
	}else{
		insert_str := strings.Join(columns,",")

		insert_sql := fmt.Sprintf("%s into %s.%s(%s)values",write_mode,to_db,to_table,insert_str)
		question_sign := make([]string,len(columns))
		for i,_ := range columns{
			question_sign[i] = "?"
		}
		question_sign_strs := fmt.Sprintf("(%s)",strings.Join(question_sign,","))
		temp_batchs := make([]map[string]string,0)
		for len(datas)>0{
			data := datas[0]
			temp_batchs = append(temp_batchs,data)
			datas = append(datas[:0],datas[1:]...)
			if len(temp_batchs) == write_batch{
				_insert_sql := insert_sql
				s := make([]string,write_batch)
				values := make([]interface{},len(columns)*len(temp_batchs))
				for index,data := range temp_batchs{
					v := make([]interface{},len(columns))
					for j :=0;j<len(columns);j++{
						if data[columns[j]] == "NULL"{
							v[j] = nil
						}else{
							v[j] = data[columns[j]]
						}
					}
					for x :=0;x<len(columns);x++{
						values[index*len(columns)+x] = v[x]
					}
					s[index] = question_sign_strs
				}
				_insert_sql += strings.Join(s,",")
				_num,err := m.Execute(_insert_sql,values...)
				if err != nil{
					status = 1
				}
				num += _num
				//fmt.Println("write ",num," rows")
				temp_batchs = make([]map[string]string,0)
			} 
		}
		if len(temp_batchs)>0{
			_insert_sql := insert_sql
			s := make([]string,len(temp_batchs))
			values := make([]interface{},len(columns)*len(temp_batchs))
			for index,data := range temp_batchs{
				v := make([]interface{},len(columns))
				for j :=0;j<len(columns);j++{
					if data[columns[j]] == "NULL"{
						v[j] = nil
					}else{
						v[j] = data[columns[j]]
					}
				}
				for x :=0;x<len(columns);x++{
					values[index*len(columns)+x] = v[x]
				}
				s[index] = question_sign_strs
			}
			_insert_sql += strings.Join(s,",")
			_num,err := m.Execute(_insert_sql,values...)
			if err != nil{
				status = 1
			}
			num += _num
			//fmt.Println("write ",num," rows")
			temp_batchs = nil
		}
		return num,is_create_table,status
    }
}
//封装query方法
func (m *MysqlClient) Query(query string,args ...interface{}) ([]map[string]string,[]string,error){
	//stmtIns, err := m.Db.Prepare(query)
	rows, err := m.Db.Query(query,args...)
	//defer stmtIns.Close()
	if err != nil {
		return nil,nil,err
	}
	columns, _ := rows.Columns()
	for index,col := range columns{
		columns[index] = strings.ToLower(col)
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]sql.RawBytes, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	results := make([]map[string]string,0)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			record[strings.ToLower(columns[i])] = value
		}
		
		results = append(results,record)
	}
	rows.Close()
	return results,columns,nil
}

func (m *MysqlClient) Execute(stmt string, args ...interface{}) (int64, error){
	result,err := m.Db.Exec(stmt,args...)
	if err != nil{
		fmt.Println("insert err ",err)
		return 0,err
	}
	var affNum int64
	 affNum, err = result.RowsAffected()
	 if err != nil {
		 fmt.Println("insert err2 ",err)
		 return 0,err
	}
	return affNum,nil
		
}


