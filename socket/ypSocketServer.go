package socket

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	//"time"
)

type SocketServer struct {
	AddressIp string
	Port int
	StatusServer bool
	ch_socketout chan SocketMessage
}

type SocketMessage struct {
	Head uint8
	Message string
}

/*******************************************************************************
*   @example
*	ch_socket := make(chan SocketMessage, 10)
*	ServerStart(10000,ch_socket)
*******************************************************************************/
func ServerStart(port int, ch_out chan SocketMessage) {
	listen, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("listen failed, err:%v\n", err)
		return
	}
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}

		go process(conn, ch_out)
	}
}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func process(conn net.Conn, ch_out chan SocketMessage) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}
		if n >0 {
			recv := recv(buf[:])
			switch recv.Head {
				case 0: {
					messageSendBack := SocketMessage{0,"Server check..."}
					send(conn, messageSendBack)
				}
			}
			ch_out <- recv
		}
	}
}

func send(conn net.Conn, message SocketMessage) bool {
	mes := []byte{message.Head}
	mes = append(mes, []byte(message.Message)...)
	_, err :=conn.Write(mes)
	if err != nil {
		return false
	} 
	return true
}

func recv(buf []byte) SocketMessage {
	var res SocketMessage
	res.Head = buf[0]
	res.Message = string(buf[1:])
	//fmt.Println(res)
	return res
}

func processDemo(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}

		recv := string(buf[:n])
		fmt.Printf("收到的数据：%v\n", recv)

		// 将接受到的数据返回给客户端
		_, err = conn.Write([]byte("ok"))
		if err != nil {
			fmt.Printf("write from conn failed, err:%v\n", err)
			break
		}
	}
}
