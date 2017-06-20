package main

import (
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	id        bson.ObjectId
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

func (u *user) castDbuser() dbuser {
	return dbuser{
		Id:        u.id,
		Name:      u.name,
		Email:     u.email,
		Password:  u.password,
		AvatarURL: u.avatarURL,
	}
}
