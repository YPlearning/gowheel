package main

import (
	//"fmt"
	_ "gowheel/sqlite"
	"gowheel/mqtt"
	"time"
)

func main(){
	var my mqtt.MQTTClient
	my.Connect("yp")
	for{
        time.Sleep(5 * time.Second)
		my.Subscribe("test",0)
		my.Publish("Topic",0,"hello")
    }
}