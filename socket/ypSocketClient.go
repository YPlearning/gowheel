package socket

import (
	"fmt"
	"net"
	"strconv"
	//"strings"
	"time"
)

type SocketClient struct {
	ServerIp string
	ServerPort int
	ClientId uint8
	StatusClient bool
	Conn *net.Conn
}

func (client *SocketClient)Start(serverIp string, port int) {
	client.ServerIp = serverIp
	client.ServerPort = port
	conn, err := net.Dial("tcp", client.ServerIp+":"+strconv.Itoa(client.ServerPort))
	if err != nil {
		fmt.Printf("conn server failed, err:%v\n", err)
		return
	}
	defer conn.Close()
	client.Conn = &conn
	client.StatusClient = true
	client.getClientId()

	//定时心跳
	heartBeatTime := 20	
	go client.heartBeat()
	go timeTicker(&heartBeatTime)

	for heartBeatTime>0 && client.StatusClient {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed:%v\n", err)
			client.StatusClient = false
			return
		}
		if n >0 {
			recv := recv(buf[:])
			switch recv.Head {
				case 0: {
					if recv.Message=="Server check..." {
						heartBeatTime = 20
					}
				}
				case 1: {
					client.ClientId = uint8(recv.Message[0])
				}
			}
		}
	}
	client.StatusClient = false


}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func (client *SocketClient)heartBeat(){
	for {
		time.Sleep(1 * time.Second)
		heartMessage := SocketMessage{0,"Client check...","END..."}
		send(*client.Conn,heartMessage)
	}
}

func (client *SocketClient)getClientId(){
	if client.StatusClient {
		heartMessage := SocketMessage{1,"Client ID...","END..."}
		send(*client.Conn,heartMessage)
	}
}

func (client *SocketClient)SendMessage(message string) {
	if message=="" {
		message = " "
	}
	if client.StatusClient {
		heartMessage := SocketMessage{10,message,"END..."}
		send(*client.Conn,heartMessage)
	}
}

func timeTicker(num *int){
	for {
		time.Sleep(1 * time.Second)
		*num--
	}
}

