package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type Message struct {
	Name string      `json: "name"`
	Data interface{} `json: "data"`
}

type Channel struct {
	Id   string `json: "id"`
	Name string `json: "name"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	// CheckOrigin: func(r *http.Request) bool {return true}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err, " 1")
		return
	}
	for {
		// msgType, msg, err := socket.ReadMessage()
		// if err != nil {
		// 	fmt.Println(err, " 2")
		// 	return
		// }
		var inMessage Message
		var outMessage Message
		if err := socket.ReadJSON(&inMessage); err != nil {
			fmt.Println(err, " 2")
			break
		}
		fmt.Printf("%#v\n", inMessage)
		switch inMessage.Name {
		case "channel add":
			err := addChannel(inMessage.Data)
			if err != nil {
				outMessage = Message{"error", err}
				if err := socket.WriteJSON(outMessage); err != nil {
					fmt.Println(err)
					break
				}
			}
			case
		}
		// fmt.Println(string(msg))
		// //Another way to handle errors
		// if err = socket.WriteMessage(msgType, msg); err != nil {
		// 	fmt.Println(err, " 3")
		// 	return
		// }

	}
}

func addChannel(data interface{}) error {
	var channel Channel
	//Type Assertion
	// channelMap := data.(map[string]interface{})
	// channel.Name = channelMap["name"].(string)
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		return err
	}
	channel.Id = "1"
	fmt.Printf("%#v\n", channel)
	return nil
}
