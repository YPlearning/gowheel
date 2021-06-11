package socket

import (
	"fmt"
	"testing"
	"time"
)

func TestSocket (t *testing.T) {
	var server SocketServer
	go server.Start(10000)

	/*
	ch_socket := make(chan SocketMessage, 10)
	go ServerStart(10000,ch_socket)*/
	
	time.Sleep(1 * time.Second)
	
	var client SocketClient
	go client.Start("127.0.0.1",10000)
	time.Sleep(1 * time.Second)
	fmt.Println(client)

	for {
		smes, ok := <- server.ch_socketout
		if ok {
			fmt.Println(smes)
		}
	}
}