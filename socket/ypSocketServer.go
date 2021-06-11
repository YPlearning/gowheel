package socket

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	//"time"
)

type SocketServer struct {
	Port int
	StatusServer bool
	ch_socketout chan SocketMessage
	ClientList [5]socketClients
}

type socketClients struct {
	ClientId uint8
	StatusClient bool
	Conn *net.Conn
}

/*******************************************************************************
*   @code
*	0	心跳检测
*	1	获取客户端ID
*	10	正常发送消息
*******************************************************************************/
type SocketMessage struct {
	Head uint8
	Message string
	End string 
}

func (server *SocketServer)Start(port int) {
	server.StatusServer = false
	server.Port = port

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}
	server.StatusServer = true
	server.ch_socketout = make(chan SocketMessage, 10)
	//打开并初始化Server完成
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		
		i, ok := server.clientConnect(conn)
		if ok {
			go server.serverListen(i)
		}
	}
}

func (server *SocketServer) clientConnect(conn net.Conn) (uint8, bool) {
	for i := range server.ClientList {
		if server.ClientList[i].StatusClient==false {
			server.ClientList[i] = socketClients{uint8(i),true,&conn}
			fmt.Println(server.ClientList)
			return uint8(i),true
		}
	}
	fmt.Println("连接数量已达上线！")
	conn.Close()
	return 0,false
}

func (server *SocketServer)serverListen(clientId uint8){
	for {
		reader := bufio.NewReader(*server.ClientList[clientId].Conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			server.ClientList[clientId].StatusClient = false
			break
		}
		//获取消息
		
		if n >0 {
			recv := recv(buf[:])
			switch recv.Head {
				case 0: {
					messageSendBack := SocketMessage{0,"Server check...","END..."}
					send(*server.ClientList[clientId].Conn, messageSendBack)
				} 
				case 1: {
					messageSendBack := SocketMessage{1,string(clientId),"END..."}
					send(*server.ClientList[clientId].Conn, messageSendBack)
				}
			}
			server.ch_socketout <- recv
		}
	}
}

func (server *SocketServer)SendToClient(clientId int,code uint8, message string) {
	if server.ClientList[clientId].StatusClient {
		mes := SocketMessage{code,message,"END..."}
		send(*server.ClientList[clientId].Conn,mes)
	}
}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func send(conn net.Conn, message SocketMessage) bool {
	mes := []byte{message.Head}
	mes = append(mes, []byte(message.Message)...)
	mes = append(mes, []byte(message.End)...)
	_, err :=conn.Write(mes)
	if err != nil {
		return false
	} 
	return true
}

func recv(buf []byte) SocketMessage {
	var res SocketMessage
	res.Head = buf[0]
	tempStr := string(buf[1:])
	n := strings.LastIndex(tempStr,"END...")
	if n==-1 {
		res.Message = ""
		res.End = "END..."
		return res
	}
	res.Message = string(buf[1:n+1])
	res.End = string(buf[n+1:n+7])
	//fmt.Println(res)
	return res
}
