package mysql

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strings"
	"strconv"
)

type MysqlClient struct {
	db *sql.DB
	IsConnect bool
}

/*******************************************************************************
*   @example
*   ch_mysql := make(chan string, 10)
*   go client.Listen(ch_mysql)
*   ch_mysql <- "tablename::'value1','value2'..."
*******************************************************************************/
func (client *MysqlClient) Listen(ch_mysql chan string) {
	for true {
		str, ok := <- ch_mysql
		if ok {
			mes := strings.Split(str,"::")
            client.Insert(mes[0], []string{mes[1]})
		}
	}
}

/*******************************************************************************
*   @example
*   var client sqlite.SqliteClient
*   client.Open("Testdb","sql123456789Server","yplearning.cn",3306,"testdb")
*******************************************************************************/
func (client *MysqlClient) Open(Username string, Password string, AddressIP string, Port int, DatabaseName string) {
	var err error
	path := Username+":"+Password+"@tcp("+AddressIP+":"+strconv.Itoa(Port)+")/"+DatabaseName+"?charset=utf8"
	//path := "Testdb:sql123456789Server@tcp(yplearning.cn:3306)/testdb?charset=utf8"
	client.db, err = sql.Open("mysql", path)
    checkErr(err)
    client.IsConnect = true
}

/*******************************************************************************
*   @example
*   client.Close()
*******************************************************************************/
func (client *MysqlClient) Close() {
	client.db.Close()
    client.IsConnect = false
}

/*******************************************************************************
*   @example
*   cloumns := []string{"ID INT PRIMARY KEY NOT NULL","name TEXT NOT NULL"}
*   client.CreateTable("hello",cloumns)
*******************************************************************************/
func (client *MysqlClient) CreateTable(tablename string, columns []string){
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
func (client *MysqlClient) DeleteTable(tablename string) {
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
*   client.Insert("hello", values)
*******************************************************************************/
func (client *MysqlClient) Insert(tablename string, values []string) int64 {
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
*   aa := client.Select("hello", "*", "ID = 111")
*   fmt.Println(aa)
*******************************************************************************/
func (client *MysqlClient) Select(tablename string, selectColumns string, condition string) []map[string]string {
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
*   aa := client.Delete("hello", "ID=111")
*   fmt.Println(aa)
*******************************************************************************/
func (client *MysqlClient) Delete(tablename string, condition string) int64 {
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

