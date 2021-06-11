package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gowheel/socket"
	"math/rand"
	"strings"
	"time"
)

type Node struct{
	NodeId int `json:"NodeId"`
	IsMaster bool `json:"IsMaster"`
	AddressIp string `json:"AddressIp"`
	Port int `json:"Port"`
}

var NodeList map[int]Node

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	myid := rand.Intn(99999999)
	myport := flag.Int("p", 10000, "ip address to run this node on. default is 10000.")
	isMaster := flag.Bool("m", false, "如果IP地址没有连接到集群中，我们将其作为Master节点.")
	flag.Parse()

	myNode := Node{myid,*isMaster,"127.0.0.1",*myport}
	NodeList = make(map[int]Node)

	if *isMaster {

	}

	ch_socket := make(chan socket.SocketMessage, 10)
	go myNode.NodeStart(ch_socket)

	for true {
		mes, ok := <- ch_socket
		if ok {
			switch mes.Head {
				case 1:	{
					nodeAdd(mes.Message)
					fmt.Println(NodeList)
				}
			}
			fmt.Println(mes)
		}
	}

	fmt.Println(myNode)
}

/*******************************************************************************
*	@instructions
*	{1 "Node"}
*******************************************************************************/
func (node *Node)NodeStart(ch_socket chan socket.SocketMessage){
	//启动Socket Server和Client
	var server socket.SocketServer
	go server.Start(10000)
	time.Sleep(1 * time.Second)
	var client socket.SocketClient
	go client.Start("127.0.0.1",10000)


}

/*******************************************************************************
*   Internal Function
*******************************************************************************/
func nodeAdd(mes string){
	var tempNode Node
	idx := strings.Index(mes,"}")
    json.Unmarshal([]byte(mes[:idx+1]), &tempNode)
	NodeList[tempNode.NodeId] = tempNode
}