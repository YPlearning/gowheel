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
*   @example
*   ch_sqlite := make(chan string, 10)
*   go client.Listen(ch_sqlite)
*   ch_sqlite <- "tablename::'value1','value2'..."
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
*   @example
*   var my sqlite.SqliteClient
*   client.Open("./test.db")
*******************************************************************************/
func (client *SqliteClient) Open(filePath string) {
	var err error
	client.db, err = sql.Open("sqlite3", filePath)
    checkErr(err)
    client.IsConnect = true
}

/*******************************************************************************
*   @example
*   client.Close()
*******************************************************************************/
func (client *SqliteClient) Close(filePath string) {
	client.db.Close()
    client.IsConnect = false
}

/*******************************************************************************
*   @example
*   cloumns := []string{"ID INT PRIMARY KEY NOT NULL","name TEXT NOT NULL"}
*   client.CreateTable("test",cloumns)
*******************************************************************************/
func (client *SqliteClient) CreateTable(tablename string, columns []string){
    if !client.IsConnect {
        fmt.Println("Database Disconnected!!!")
        return
    }
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
*   @example
*   client.DeleteTable("hello")
*******************************************************************************/
func (client *SqliteClient) DeleteTable(tablename string) {
    if !client.IsConnect {
        fmt.Println("Database Disconnected!!!")
        return
    }
    //sql := "DROP TABLE test"
    sqlStr := "DROP TABLE " + tablename
    _, err := client.db.Exec(sqlStr)
	checkErr(err)
}

/*******************************************************************************
*   @example
*   values := []string{"'111'","'aaa'"}
*   client.Insert("test", values)
*******************************************************************************/
func (client *SqliteClient) Insert(tablename string, values []string) int64 {
    if !client.IsConnect {
        fmt.Println("Database Disconnected!!!")
        return 0
    }
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
*   @example
*   aa := client.Select("test", "*", "ID = 111")
*   fmt.Println(aa)
*******************************************************************************/
func (client *SqliteClient) Select(tablename string, selectColumns string, condition string) []map[string]string {
    if !client.IsConnect {
        fmt.Println("Database Disconnected!!!")
        return nil
    }
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
*   @example
*   aa := client.Delete("test", "ID=111")
*   fmt.Println(aa)
*******************************************************************************/
func (client *SqliteClient) Delete(tablename string, condition string) int64 {
    if !client.IsConnect {
        fmt.Println("Database Disconnected!!!")
        return 0
    }
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

