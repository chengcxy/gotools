package backends

import (
	"fmt"
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
	mc.ConnUri = mc.getConnUri()
	return mc
	    

}

//拼接url
func (mc *MysqlConfig) getConnUri()(url string){
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








//mysql客户端结构体
type MysqlClient struct{
	Config *MysqlConfig
	Db *sql.DB
}

//mysql 客户端
func NewMysqlClient(mc *MysqlConfig)(*MysqlClient){
	Db,_ :=  sql.Open("mysql",mc.ConnUri)
	Db.SetConnMaxLifetime(time.Minute * 100)
	Db.SetMaxOpenConns(20)
	Db.SetMaxIdleConns(20)
	client := &MysqlClient{
		Config:mc,
		Db:Db,
	}
	return client
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
	//stmtIns, _ := m.Db.Prepare(stmt)
	result,err := m.Db.Exec(stmt,args...)
	//defer stmtIns.Close()
	rows_affected ,_ := result.RowsAffected()
	return rows_affected,err
}