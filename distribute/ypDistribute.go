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
	AddressIp string `json:"AddressIp"`
	Port int `json:"Port"`
}

var NodeList map[int]Node

func main(){
	rand.Seed(time.Now().UTC().UnixNano())
	myid := rand.Intn(99999999)
	myport := flag.Int("p", 10000, "ip address to run this node on. default is 10000.")
	flag.Parse()

	myNode := Node{myid,"127.0.0.1",*myport}
	NodeList = make(map[int]Node)

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
	go socket.ServerStart(10000,ch_socket)
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