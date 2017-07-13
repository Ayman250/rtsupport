package main

import (
"github.com/mitchellh/mapstructure"
r "github.com/dancannon/gorethink"
)

func addChannel(client *Client, data interface{}) {
	var channel Channel
	err := mapstructure.Decode(data, &channel)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}
	/*Wrapping db access in self invoked anyonymous function that will run in it's own go routine. 
	This is because it's a slow blocking process. No reason have it running in main thread*/
	go func() {
		err := r.Table("channel").Insert(channel).Exec(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}