package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	// "encoding/json"
)

var testMsg Message

type FindHandler func(string) (Handler, bool)

type Message struct {
	Name string      `json: "name"`
	Data interface{} `json: "data"`
}

type Client struct{
	send chan Message
	socket *websocket.Conn
	findHandler FindHandler
}


func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			fmt.Println(err)
			break
		}
	}

	if handler, found := client.findHandler(message.Name); found {
		handler(client, message.Data)
	}
	client.socket.Close()
}


func (client *Client) Write(){
	fmt.Println("WRITE")
	if testMsg.Name != "" {
		client.socket.WriteJSON(testMsg)
	} else {
		fmt.Println("ITS NILLL")
	}
	for msg := range client.send {
		testMsg = msg
		// if err := client.socket.WriteJSON(msg); err != nil {
		// 	fmt.Println(err)
		// 	break
		// }
	}
	client.socket.Close();
}


func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client{
	return &Client{
		send: make(chan Message),
		socket: socket,
		findHandler: findHandler,
	}
}



