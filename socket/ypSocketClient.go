package socket

import (
	"fmt"
	"net"
	"strconv"
	//"strings"
	"time"
)

/*******************************************************************************
*   @example
*	ch_socketin := make(chan SocketMessage, 10)
*	go ClientStart("127.0.0.1", 10000, ch_socketin)
*******************************************************************************/
func ClientStart(address string, port int, ch_in chan SocketMessage){
	conn, err := net.Dial("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return
	}
	defer conn.Close()

	//心跳检测
	go heartBeat(conn)
	go listen(conn,ch_in)
	heartBeatTime := 20
	go timeTicker(&heartBeatTime)
	for heartBeatTime>0 {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed:%v\n", err)
			return
		}
		if n >0 {
			recv := recv(buf[:16])
			switch recv.Head {
				case 0: {
					if recv.Message=="Server check..." {
						heartBeatTime = 20
					}
				}
			}
		}
	}
}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func heartBeat(conn net.Conn){
	for {
		time.Sleep(1 * time.Second)
		heartMessage := SocketMessage{0,"Client check..."}
		send(conn,heartMessage)
	}
}

func timeTicker(num *int){
	for {
		time.Sleep(1 * time.Second)
		*num--
	}
}

func listen(conn net.Conn, ch_in chan SocketMessage){
	smes, ok := <- ch_in
	if ok {
		send(conn,smes)
	}
}