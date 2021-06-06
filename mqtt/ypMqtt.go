package mqtt

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
	"fmt"
)

type MQTTClient struct {
	client mqtt.Client
	IsConnect bool
} 

/*******************************************************************************
//@example
//var client mqtt.MQTTClient
//client.Connect("testID")
*******************************************************************************/
func (client *MQTTClient) Connect(clientID string) {
	var url = "tcp://192.168.1.108:1883"
    opts := mqtt.NewClientOptions()
    opts.AddBroker(url)//fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID(clientID)
    //opts.SetUsername("nameYP")
    //opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client.client = mqtt.NewClient(opts)
    if token := client.client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
  	} else {
		client.IsConnect = true
	}
}

/*******************************************************************************
//@example
//client.Publish("Topic",0,"hello")
*******************************************************************************/
func (client *MQTTClient) Publish(topic string, qos byte, payload interface{}) {
	client.client.Publish(topic, qos, false, payload)
}

/*******************************************************************************
//@example
//client.Subscribe("testTopic",0)
*******************************************************************************/
func (client *MQTTClient) Subscribe(topic string, qos byte){
	client.client.Subscribe(topic,qos,generalCallback)
}

/*******************************************************************************
*    Callback Function
*******************************************************************************/
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
	client.Subscribe("Topic",0x00,topicCallback)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func topicCallback(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Subscribe: Topic is [%s]; msg is [%s]\n", msg.Topic(), string(msg.Payload()))
}

func generalCallback(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("General Subscribe: Topic is [%s]; msg is [%s]\n", msg.Topic(), string(msg.Payload()))
}
