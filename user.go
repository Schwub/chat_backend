package main

import ()

type user struct {
	id        uint64
	name      string
	email     string
	password  string
	avatarURL string
}

func (u user) userJson() map[string]interface{} {
	uMap := make(map[string]interface{})
	uMap["id"] = u.id
	uMap["name"] = u.name
	uMap["avatarURL"] = u.avatarURL
	uMap["email"] = u.email
	return uMap
}
