package main

import (
	"log"
)

func newMessage(c *client, m map[string]interface{}) {
	d := m["data"]
	data := d.(map[string]interface{})
	newmessage := make(map[string]interface{})
	newmessage["type"] = "event"
	newmessage["subtype"] = "chat"
	newmessage["event"] = "newMessage"
	messageData := make(map[string]interface{})
	messageData["message"] = data["message"]
	messageData["roomName"] = data["channelName"].(string)
	newmessage["data"] = messageData
	log.Println("----------------------------------------------", newmessage)
	c.hub.rooms[data["channelName"].(string)].sendToRoom(newmessage)

}
