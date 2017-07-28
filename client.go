package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
)

type FindHandler func(string) (Handler, bool)



type Client struct{
	send chan Message
	socket *websocket.Conn
	findHandler FindHandler
	session *r.Session
	stopChannels map[int]chan bool
	id string
	userName string
}

func (c *Client) NewStopChannel(stopKey int) chan bool{
	/*This is a safety mechanism incase client calls channel subscribe multiple times without unsubscribe
	by deleteing the old channel before adding it, it prevents go routine leakes*/
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func(c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found {
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
	client.socket.Close()
}

func (c *Client) Close(session *r.Session){
	for _, ch := range c.stopChannels {
		ch <- true
	}
	/*After disconnecting, delete user form database based on user name (Bad filter, should use key instead)*/
	r.Table("user").Filter(map[string]interface{}{
    		"name": c.userName,
		}).Delete().RunWrite(session)
	fmt.Println(c.userName)
	/*Closing the send channel ends the for loop in the write effectively killing that go routine
	Clever way to clean up instead of using another stop channel for the write*/
		close(c.send)
}


func NewClient(socket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client{
	var user User
	user.Name = "Anonymous"
	res, err := r.Table("user").Insert(user).RunWrite(session)
	if err != nil {
		log.Println(err.Error())
	}
	var id string
	if len(res.GeneratedKeys) > 0 {
		id = res.GeneratedKeys[0]
	}
	return &Client{
		send: make(chan Message),
		socket: socket,
		findHandler: findHandler,
		session: session,
		stopChannels: make(map[int]chan bool),
		id: id,
		userName: "Anonymous",
	}
}