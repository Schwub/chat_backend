package main

import ()

func getAllUsers(h *hub) interface{} {

	//users := make([]user, 0, len(m))
	//for k := range m {
	//	users = append(users, k.user)
	//}
	//allusers := allUsers{
	//	mtype:   "event",
	//	subtype: "user",
	//	event:   "allUsers",
	//	data:    users,
	//}
	alluser := make(map[string]interface{})
	alluser["type"] = "event"
	alluser["subtype"] = "user"
	alluser["event"] = "allUsers"
	alluser["data"] = h.usersJson()
	return alluser
	//rainer := make(map[string]string)
	//rainer["name"] = "Rainer Winkler"
	//rainer["email"] = "drache@offiziel.alt"
	//rainer["id"] = "eins"

	//users := make([]interface{}, 0, 1)
	//users = append(users, rainer)
	//alluser["data"] = users
	//return alluser
}

func newUserEvent(u user) interface{} {
	newuser := make(map[string]interface{})
	newuser["type"] = "event"
	newuser["subtype"] = "user"
	newuser["event"] = "newUser"
	newuser["data"] = u.userJson()
	return newuser
}

func userLeaves(c *client) interface{} {
	userleaves := make(map[string]interface{})
	userleaves["type"] = "event"
	userleaves["subtype"] = "user"
	userleaves["event"] = "userLeaves"
	userleaves["data"] = c.user.userJson()
	return userleaves
}
