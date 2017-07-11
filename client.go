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
}

func (c *client) read() {
	for {
		var dat map[string]interface{}
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			json.Unmarshal(msg, &dat)
			log.Printf(c.user.name, ": client.read()", dat)
			c.handleMessage(dat)
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		log.Println(c.user.name, ": client.write() ", msg)
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (c *client) handleMessage(m map[string]interface{}) {
	log.Println("handle Message")
	switch m["subtype"] {
	case "user":
		c.handleUserEvent(m)
	case "auth":
		c.handleAuthEvent(m)
	case "room":
		c.handleRoomEvent(m)
	case "message":
		c.handleMessageEvent(m)
	default:
		log.Println("Message handling not implemented, yet")
	}
}

func (c *client) handleUserEvent(m map[string]interface{}) {
	switch m["command"] {
	case "getAllUsers":
		log.Println("handle getAllUsers")
		msg := getAllUsers(c.hub)
		c.send <- msg
	}
}

func (c *client) handleAuthEvent(m map[string]interface{}) {
	switch m["command"] {
	case "newRegistration":
		log.Println("handle newRegistration")
		msg := newRegistration(c, m)
		c.send <- msg
		msg = newUserEvent(c.user)
		c.hub.sendToAll(msg)
		c.hub.db.DB("chat").C("users").Insert(c.user.castDbuser())
	case "logout":
		log.Println("handle logout")
		logout(c, m)
		msg := userLeaves(c)
		c.hub.sendToAll(msg)
	case "login":
		log.Println("handle login")
		msg := login(c, m)
		c.send <- msg
		msg = newUserEvent(c.user)
		c.hub.sendToAll(msg)
	case "checkLogin":
		log.Println("handle checkLogin")
		msg := checkLogin(c, m)
		c.send <- msg
		msg = newUserEvent(c.user)
		c.hub.sendToAll(msg)
	}
}

func (c *client) handleRoomEvent(m map[string]interface{}) {
	switch m["command"] {
	case "createRoom":
		log.Println("handle newRoom")
		msg, room := createRoom(c, m)
		c.send <- msg
		msg = newRoomEvent(*room)
		c.hub.sendToAll(msg)

	case "joinRoom":
		log.Println("handle joinRoom")
		msg := joinRoom(c, m)
		c.send <- msg
	case "leaveRoom":
		log.Println("handle leaveRoom")
		leaveRoom(c, m)
	case "getAllRooms":
		//TODO
		log.Println("handle getAllRooms")
		msg := getAllRooms(c, m)
		c.send <- msg
	}
}

func (c *client) handleMessageEvent(m map[string]interface{}) {
	switch m["command"] {
	case "newMessage":
		log.Println("handle newMessage")
		newMessage(c, m)
	}
}
