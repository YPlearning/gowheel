package mysql

import (
	"fmt"
	"testing"
	"time"
)

func TestMysql (t *testing.T){
	fmt.Println("hello")
	var client MysqlClient
    client.Open("Testdb","sql123456789Server","yplearning.cn",3306,"testdb")

	cloumns := []string{"ID INT PRIMARY KEY NOT NULL","name TEXT NOT NULL"}
    client.CreateTable("hello",cloumns)

	values := []string{"'111'","'aaa'"}
    client.Insert("hello", values)

	ch_mysql := make(chan string, 10)
    go client.Listen(ch_mysql)
    ch_mysql <- "hello::'222','ch_mysql'"
	time.Sleep(1 * time.Second)

	aa := client.Select("hello", "*", "ID = 111")
    fmt.Println(aa)

	bb := client.Delete("hello", "ID")
    fmt.Println(bb)

	client.DeleteTable("hello")

	client.Close()
	fmt.Println("yes")
}