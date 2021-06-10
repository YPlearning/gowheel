package socket

import (
	"fmt"
	"testing"
	"time"
)

func TestSocket (t *testing.T) {

	ch_socket := make(chan SocketMessage, 10)
	go ServerStart(10000,ch_socket)
	time.Sleep(1 * time.Second)
	
	ch_socketin := make(chan SocketMessage, 10)
	go ClientStart("127.0.0.1", 10000, ch_socketin)
	
	mes := SocketMessage{1,"hello"}
	ch_socketin <- mes

	for true {
		smes, ok := <- ch_socket
		if ok {
			fmt.Println(smes)
		}
	}
}