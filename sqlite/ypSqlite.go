package sqlite

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strings"
)

type SqliteClient struct {
	db *sql.DB
	IsConnect bool
}

/*******************************************************************************
//@example
//ch_sqlite := make(chan string, 10)
//go client.Listen(ch_sqlite)
//ch_sqlite <- "tablename::'value1','value2'..."
*******************************************************************************/
func (client *SqliteClient) Listen(ch_sqlite chan string) {
	for true {
		str, ok := <- ch_sqlite
		if ok {
			mes := strings.Split(str,"::")
            client.Insert(mes[0], []string{mes[1]})
		}
	}
}

/*******************************************************************************
//@example
//var my sqlite.SqliteClient
//client.Open("./test.db")
*******************************************************************************/
func (client *SqliteClient) Open(filePath string) {
	var err error
	client.db, err = sql.Open("sqlite3", filePath)
    checkErr(err)
    client.IsConnect = true
}

/*******************************************************************************
//@example
//client.Close()
*******************************************************************************/
func (client *SqliteClient) Close(filePath string) {
	client.db.Close()
    client.IsConnect = false
}

/*******************************************************************************
//@example
//cloumns := []string{"ID INT PRIMARY KEY NOT NULL","name TEXT NOT NULL"}
//client.CreateTable("test",cloumns)
*******************************************************************************/
func (client *SqliteClient) CreateTable(tablename string, columns []string){
	n := len(columns)
	//sql := "CREATE TABLE test (ID INT PRIMARY KEY NOT NULL,NAME TEXT NOT NULL)" 
	sqlStr := "CREATE TABLE " + tablename + " ("
	for i := 0; i < n; i++ {
		sqlStr += columns[i]
		if i<n-1 {
			sqlStr += ","
		}
	}
	sqlStr += ")"
	_, err := client.db.Exec(sqlStr)
	checkErr(err)
}

/*******************************************************************************
//@example
//client.DeleteTable("hello")
*******************************************************************************/
func (client *SqliteClient) DeleteTable(tablename string) {
    //sql := "DROP TABLE test"
    sqlStr := "DROP TABLE " + tablename
    _, err := client.db.Exec(sqlStr)
	checkErr(err)
}

/*******************************************************************************
//@example
//values := []string{"'111'","'aaa'"}
//client.Insert("test", values)
*******************************************************************************/
func (client *SqliteClient) Insert(tablename string, values []string) int64 {
    //sql := "INSERT INTO test values(?,?)"
    n := len(values)
    sqlStr := "INSERT INTO " + tablename + " values("
    for i := 0; i < n; i++ {
        sqlStr +=  values[i] 
        if i<n-1 {
            sqlStr += ","
        }
    }
    sqlStr += ")"

    res, err := client.db.Exec(sqlStr)
    checkErr(err)
    if res==nil {
        return 0
    }
    id, err := res.RowsAffected()
    checkErr(err)

    return id
}

/*******************************************************************************
//@example
//aa := client.Select("test", "*", "ID = 111")
//fmt.Println(aa)
*******************************************************************************/
func (client *SqliteClient) Select(tablename string, selectColumns string, condition string) []map[string]string {
    sqlStr := "SELECT " + selectColumns + " FROM " + tablename 
    if condition!="" {
        sqlStr += " WHERE " + condition
    }
	rows, err := client.db.Query(sqlStr)
    checkErr(err)
    if rows==nil {
        return nil
    }
    columns, err := rows.Columns()
    checkErr(err)

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
        scanArgs[i] = &values[i]
    }
 	ccc,err := rows.Columns()
	var slice []map[string]string
	for rows.Next(){
		var myMap map[string]string
        myMap = make(map[string]string)
		rows.Scan(scanArgs...)
		var value string
        for i , col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
				myMap[ccc[i]] = value
            }
		}
		slice = append(slice, myMap)
	}
	rows.Close()
	return slice
}

/*******************************************************************************
@example
//aa := client.Delete("test", "ID=111")
//fmt.Println(aa)
*******************************************************************************/
func (client *SqliteClient) Delete(tablename string, condition string) int64 {
    sqlStr := "delete from " + tablename + " where "
    if condition=="" {
        return 0
    } else {
        sqlStr += condition
    }
    res, err := client.db.Exec(sqlStr)
    checkErr(err)
    if res==nil {
        return 0
    }

    id, err := res.RowsAffected()
    checkErr(err)

    return id
}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func checkErr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}




/*
func Test() {
    db, err := sql.Open("sqlite3", "./foo.db")
    checkErr(err)
	
	//新建数据表
	sql := "CREATE TABLE test (ID INT PRIMARY KEY NOT NULL,NAME TEXT NOT NULL)" 
	res, err := db.Exec(sql)
	checkErr(err)
	fmt.Println(res)

	
	//插入数据
    stmt, err := db.Prepare("INSERT INTO test(ID, NAME) values(?,?)")
    checkErr(err)

    res, err := stmt.Exec("222", "hello")
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)

    //更新数据
    stmt, err := db.Prepare("update test set NAME=? where ID=?")
    checkErr(err)

    res, err := stmt.Exec("xiexie", "111")
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)

    //查询数据
    rows, err := db.Query("SELECT * FROM test")
    checkErr(err)

    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        checkErr(err)
        fmt.Println(id,name)
    }

    //删除数据
    stmt, err := db.Prepare("delete from test where ID=?")
    checkErr(err)

    res, err := stmt.Exec("111")
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)
/*	

    db.Close()

}*/

