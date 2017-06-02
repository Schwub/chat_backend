package main

import (
	"github.com/fatih/structs"
)

type allUsers struct {
	mtype   string `json:"type"`
	subtype string `json:subtype`
	event   string `json:event`
	data    []user `json:data`
}

func allUsersJson(m map[*client]bool) interface{} {
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
	//rainer := make(map[string]string)
	//rainer["name"] = "Rainer Winkler"
	//rainer["email"] = "drache@offiziel.alt"
	//rainer["id"] = "eins"
	rainer := user{
		name:  "Rainer Winkler",
		email: "drache@offiziel.alt",
		id:    1,
	}

	users := make([]interface{}, 0, 1)
	users = append(users, structs.Map(rainer))
	alluser["data"] = users
	return alluser
}
