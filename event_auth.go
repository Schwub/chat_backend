package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

func newRegistration(c *client, m map[string]interface{}) interface{} {
	d := m["data"]
	data := d.(map[string]interface{})
	emails := c.hub.createUserEmailMap()
	log.Println(emails)
	log.Println(emails[data["email"].(string)])
	if _, registered := emails[data["email"].(string)]; registered {
		registrationError := make(map[string]interface{})
		registrationError["type"] = "error"
		registrationError["subtype"] = "auth"
		registrationError["error"] = "registrationError"
		registrationError["data"] = "Email already registered"
		return registrationError
	}
	if _, registered := c.hub.registerdUsers[data["name"].(string)]; !registered {
		newUser := user{
			name:     data["name"].(string),
			email:    data["email"].(string),
			password: data["password"].(string),
			id:       bson.NewObjectId(),
		}
		c.hub.registerdUsers[data["name"].(string)] = &newUser
		c.user = newUser
		registrationsucces := make(map[string]interface{})
		registrationsucces["type"] = "event"
		registrationsucces["subtype"] = "auth"
		registrationsucces["event"] = "registrationSuccess"
		user := make(map[string]interface{})
		user["user"] = c.user.userJson()
		registrationsucces["data"] = user
		return registrationsucces
	}
	registrationError := make(map[string]interface{})
	registrationError["type"] = "error"
	registrationError["subtype"] = "auth"
	registrationError["error"] = "registrationError"
	registrationError["data"] = "Name already registered"
	return registrationError
}

func login(c *client, m map[string]interface{}) map[string]interface{} {
	log.Println("----------------______________--------------------", m)
	d := m["data"]
	data := d.(map[string]interface{})
	for _, user := range c.hub.registerdUsers {
		if user.email == data["user"].(string) {
			if user.password != data["password"].(string) {
				authError := make(map[string]interface{})
				authError["type"] = "error"
				authError["subtype"] = "auth"
				authError["error"] = "authError"
				authError["data"] = "Wrong Password"
				return authError
			} else {
				c.user = *user
				authsucces := make(map[string]interface{})
				authsucces["type"] = "event"
				authsucces["subtype"] = "auth"
				authsucces["event"] = "authSuccess"
				userJson := make(map[string]interface{})
				userJson["user"] = c.user.userJson()
				authsucces["data"] = userJson
				return authsucces
			}
		}
	}
	authError := make(map[string]interface{})
	authError["type"] = "error"
	authError["subtype"] = "auth"
	authError["error"] = "authError"
	authError["data"] = "User does not exist"
	return authError
}

func logout(c *client, m map[string]interface{}) {
	c = nil
}

func checkLogin(c *client, m map[string]interface{}) map[string]interface{} {
	d := m["data"]
	data := d.(map[string]interface{})
	u := data["user"]
	authuser := u.(map[string]interface{})
	if authuser["name"] != "" {
		for _, user := range c.hub.registerdUsers {
			if user.name == authuser["name"] {
				c.user = *user
				checklogin := make(map[string]interface{})
				checklogin["type"] = "event"
				checklogin["subtype"] = "auth"
				checklogin["event"] = "authSucces"
				checklogin["data"] = c.user.userJson()
				return checklogin
			}
		}
	}

	errorCheckLogin := make(map[string]interface{})
	errorCheckLogin["type"] = "error"
	errorCheckLogin["subytpye"] = "auth"
	errorCheckLogin["error"] = "authError"
	errorCheckLogin["data"] = "error"
	return errorCheckLogin
}
