package smservice

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

// UserMessage model
type UserMessage struct {
	Id         string
	UserID     string
	Type       string
	AppType    int
	Title      string
	Content    string
	SendWay    string
	ReadStatus int
	FromUid    string
	BankId     string
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

const (
	DBTYPE = "mysql"
)

func Select() {

	query, err := db.Query(`select * from users`)
	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	printResult(query)

}

// todo 通用模型
func Add(um UserMessage) int64 {

	stmt, _ := db.Prepare(`INSERT user_message (id,user_id,type,app_type,title,content,send_way,read_status,from_uid,bank_id,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?,?,?)`)

	res, _ := stmt.Exec(um.Id, um.UserID, um.Type, um.AppType, um.Title, um.Content, um.SendWay, um.ReadStatus, um.FromUid, um.BankId, um.CreatedAt, um.UpdatedAt)

	row, _ := res.RowsAffected()
	return row
}

func printResult(query *sql.Rows) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	for k, v := range results { //查询出来的数组
		fmt.Println(k, v)
	}
}
