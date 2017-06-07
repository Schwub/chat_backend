package main

import ()

type allUsers struct {
	mtype   string `json:"type"`
	subtype string `json:subtype`
	event   string `json:event`
	data    []user `json:data`
}

func allUsersJson(h hub) interface{} {

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
