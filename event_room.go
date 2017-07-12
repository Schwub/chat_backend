package main

import ()

func createRoom(c *client, m map[string]interface{}) (interface{}, *room) {

	room := newRoom(m["data"].(string))
	for _, eroom := range c.hub.rooms {
		if room.name == eroom.name {
			channelError := make(map[string]interface{})
			channelError["type"] = "error"
			channelError["subtype"] = "room"
			channelError["error"] = "newChannelError"
			channelError["data"] = "Chat Name already in use"
			return channelError, nil
		}
	}
	c.hub.rooms[room.name] = room

	createroom := make(map[string]interface{})
	createroom["type"] = "event"
	createroom["subtype"] = "room"
	//TODO handle room new already in user
	createroom["event"] = "newChannelSuccess"
	createroom["data"] = room.name
	//c.hub.rooms[room.name].clients[c] = true
	return createroom, room
}

func newRoomEvent(r room) interface{} {
	newroomevent := make(map[string]interface{})
	newroomevent["type"] = "event"
	newroomevent["subtype"] = "room"
	newroomevent["event"] = "newRoom"
	newroomevent["data"] = r.roomJson()
	return newroomevent
}

func joinRoom(c *client, m map[string]interface{}) interface{} {
	d := m["data"]
	data := d.(map[string]interface{})
	r := data["channel"]
	room := r.(map[string]interface{})
	c.hub.rooms[room["name"].(string)].clients[c] = true
	joinroom := make(map[string]interface{})
	joinroom["type"] = "event"
	joinroom["subtype"] = "room"
	joinroom["event"] = "joinRoomSuccess"
	newMember(c.user, c.hub.rooms[room["name"].(string)], c)
	//TODO Error

	return joinroom
}

func leaveRoom(c *client, m map[string]interface{}) {
	d := m["data"]
	data := d.(map[string]interface{})
	memberLeaves(c.user, data["name"].(string), c)
	delete(c.hub.rooms[data["name"].(string)].clients, c)
	if len(c.hub.rooms[data["name"].(string)].clients) == 0 {
		delete(c.hub.rooms, data["name"].(string))
		msg := getAllRooms(c, m)
		c.hub.sendToAll(msg)
	}
}

func getAllRooms(c *client, m map[string]interface{}) interface{} {
	getallrooms := make(map[string]interface{})
	getallrooms["type"] = "event"
	getallrooms["subtype"] = "room"
	getallrooms["event"] = "allRooms"
	getallrooms["data"] = c.hub.roomsJson()
	return getallrooms
}

func newMember(u user, r *room, c *client) {
	newmember := make(map[string]interface{})
	newmember["type"] = "event"
	newmember["subtype"] = "room"
	newmember["event"] = "newMember"
	data := make(map[string]interface{})
	data["name"] = r.name
	data["user"] = u.userJson()
	newmember["data"] = data
	for c := range c.hub.clients {
		c.send <- newmember
	}
}

func memberLeaves(u user, r string, c *client) {
	memberleaves := make(map[string]interface{})
	memberleaves["type"] = "event"
	memberleaves["subtype"] = "room"
	memberleaves["event"] = "leaveRoom"
	data := make(map[string]interface{})
	data["channelName"] = r
	data["userId"] = u.id
	memberleaves["data"] = data
	for c := range c.hub.clients {
		c.send <- memberleaves
	}
}
