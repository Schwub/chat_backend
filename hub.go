package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type hub struct {
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	rooms   map[room]string
}

func newHub() *hub {
	return &hub{
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		rooms:   make(map[room]string),
	}
}

func (h *hub) run() {
	log.Println("Running Hub")
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
			log.Println("Client joined")
		case client := <-h.leave:
			delete(h.clients, client)
			close(client.send)
			log.Println("Client left")
		}
	}
}

const (
	socketBuffersize  = 1024
	messageBuffersize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBuffersize,
	WriteBufferSize: socketBuffersize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *hub) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("ServeHTTP:", err)
		return
	}
	client := &client{
		hub:    h,
		socket: socket,
		send:   make(chan interface{}),
		rooms:  make(map[*room]string),
	}
	h.join <- client
	defer func() { h.leave <- client }()
	go client.write()
	client.read()
}
