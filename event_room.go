package main

import ()

func createRoom(c *client, m map[string]interface{}) (interface{}, *room) {
	room := newRoom(m["data"].(string))
	c.hub.rooms[room.name] = room
	createroom := make(map[string]interface{})
	createroom["type"] = "event"
	createroom["subtype"] = "room"
	//TODO handle room new already in user
	createroom["event"] = "newChannelSuccess"
	createroom["data"] = room.name
	c.hub.rooms[room.name].clients[c] = true
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
	//TODO Error

	return joinroom
}

func leaveRoom(c *client, m map[string]interface{}) {
	d := m["data"]
	data := d.(map[string]interface{})
	delete(c.hub.rooms[data["name"].(string)].clients, c)
}

func getAllRooms(c *client, m map[string]interface{}) interface{} {
	getallrooms := make(map[string]interface{})
	getallrooms["type"] = "event"
	getallrooms["subytype"] = "room"
	getallrooms["event"] = "allRooms"
	getallrooms["data"] = c.hub.roomsJson()
	return getallrooms
}
