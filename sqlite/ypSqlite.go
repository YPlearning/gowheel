package sqlite

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

type SqliteClient struct {
	db *sql.DB
	isConnect bool
}

func (client *SqliteClient) Open(filePath string) {
	var err error
	client.db, err = sql.Open("sqlite3", filePath)
    checkErr(err)
}

func (client *SqliteClient) Select() {
	rows, err := client.db.Query("SELECT * FROM test")
    checkErr(err)

    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        checkErr(err)
        fmt.Println(id,name)
    }
}

func (client *SqliteClient) CreateTable(tablename string, columns []string){
	n := len(columns)
	//sql := "CREATE TABLE test (ID INT PRIMARY KEY NOT NULL,NAME TEXT NOT NULL)" 
	sql := "CREATE TABLE " + tablename + " ("
	for i := 0; i < n; i++ {
		sql += columns[i]
		if i<n-1 {
			sql += ","
		}
	}
	sql += ")"
	fmt.Println(sql)
	res, err := client.db.Exec(sql)
	checkErr(err)
	fmt.Println(res)
}

func Test() {
    db, err := sql.Open("sqlite3", "./foo.db")
    checkErr(err)
	/*
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
    }*/

    //删除数据
    stmt, err := db.Prepare("delete from test where ID=?")
    checkErr(err)

    res, err := stmt.Exec("111")
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)
/*	*/

    db.Close()

}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}