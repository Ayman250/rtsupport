package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	// "encoding/json"
)

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
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			client.socket.WriteJSON(msg)
			fmt.Println(msg.Name, '\n', msg.Data)
			fmt.Println(err)
		}
	}
	client.socket.Close();
	// for msg := range client.send {
	// 	w, err := client.socket.NextWriter(websocket.TextMessage)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	err1 := json.NewEncoder(w).Encode(msg)
	// 	err2 := w.Close()
	// 	if err1 != nil {
	// 		fmt.Println(err1)
	// 	}
	// 	fmt.Println(err2)
	// }
}


func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client{
	return &Client{
		send: make(chan Message),
		socket: socket,
		findHandler: findHandler,
	}
}



