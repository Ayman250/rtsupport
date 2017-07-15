package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	r "github.com/dancannon/gorethink"
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
	session *r.Session
}


func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(message.Name)
		if handler, found := client.findHandler(message.Name); found {
			fmt.Println("Handler Found")
			handler(client, message.Data)
		}
	}

	client.socket.Close()
}


func (client *Client) Write(){
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			fmt.Println(err)
			break
		}
	}
	client.socket.Close();
}


func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client{
	return &Client{
		send: make(chan Message),
		socket: socket,
		findHandler: findHandler,
		session: session,
	}
}