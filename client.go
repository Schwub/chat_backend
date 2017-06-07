package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	hub    *hub
	socket *websocket.Conn
	user   user
	send   chan interface{}
	rooms  map[*room]string
}

func (c *client) read() {
	for {
		var dat map[string]interface{}
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			json.Unmarshal(msg, &dat)
			log.Printf("client/read()", dat)
			c.handleMessage(dat)
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		log.Println("client.write: ", msg)
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (c client) handleMessage(m map[string]interface{}) {
	log.Println("handle Message")
	switch m["subtype"] {
	case "user":
		c.handleUserEvent(m)
	default:
		log.Println("Message handling not implemented, yet")
	}

}

func (c client) handleUserEvent(m map[string]interface{}) {
	switch m["command"] {
	case "getAllUsers":
		log.Println("handle allUsers")
		msg := allUsersJson(*c.hub)
		log.Println(msg)
		c.send <- msg
	}
}
