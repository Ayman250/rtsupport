package main

import(
"github.com/gorilla/websocket"
"net/http"
"fmt"

)

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Handler func(*Client, interface{})

type Router struct {
	rules map[string]Handler
}

func NewRouter() *Router{
	router := &Router{
		rules: make(map[string]Handler),
	}
	return router
}

func (r *Router) Handle(msgName string, handler Handler){
	r.rules[msgName] = handler
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

func (rout *Router) ServeHTTP(w http.ResponseWriter, r *http.Request){
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := NewClient(socket, rout.FindHandler)
	go client.Write();
	client.Read();
}