package main

import ()

type room struct {
	name    string
	clients map[*client]bool
}

func newRoom(n string) *room {
	return &room{
		name:    n,
		clients: make(map[*client]bool),
	}
}

func (r *room) roomJson() map[string]interface{} {
	mroom := make(map[string]interface{})
	mroom["name"] = r.name
	memberl := make([]interface{}, 0, 0)
	for k := range r.clients {
		memberl = append(memberl, k.user.userJson())
	}
	mroom["members"] = memberl
	return mroom
}

func (r *room) sendToRoom(message interface{}, client *client) {
	for c := range r.clients {
		if c == client {
			continue
		}
		c.send <- message
	}
}
