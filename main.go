package main

import (
	//"fmt"
	"gowheel/sqlite"
	"gowheel/mqtt"
	"time"
	"strconv"
)

func main(){
	var my mqtt.MQTTClient
	my.Connect("yp")
	ch_mqtt := make(chan string, 10)
	go my.Listen(ch_mqtt)

	var mysql sqlite.SqliteClient
	mysql.Open("./foo.db")
	ch_sqlite := make(chan string, 10)
	go mysql.Listen(ch_sqlite)

	ch_sqlite <- "test::'111','aaa'"

	for i := 0; i < 100; i++ {
		str := strconv.Itoa(i)
		go newmqtt("yp"+str)
	}

	for true {
		time.Sleep(5 * time.Second)
	}
}

func newmqtt(id string){
	var my mqtt.MQTTClient
	my.Connect(id)
}