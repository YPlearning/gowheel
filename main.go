package main

import "gowheel/sqlite"

func main(){
	//sqlite.Test()
	var my sqlite.SqliteClient
	my.Open("./foo.db")
	my.Select()
	cloums := []string{"ID INT PRIMARY KEY NOT NULL","NAME TEXT NOT NULL"}
	my.CreateTable("test1",cloums)
}