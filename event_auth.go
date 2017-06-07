package main

import ()

func newRegistration(c *client, m map[string]interface{}) interface{} {
	d := m["data"]
	data := d.(map[string]interface{})
	if _, registered := c.hub.registerdUsers[data["name"].(string)]; !registered {
		newUser := user{
			name:     data["name"].(string),
			email:    data["email"].(string),
			password: data["password"].(string),
		}
		c.hub.registerdUsers[data["name"].(string)] = &newUser
		c.user = newUser
		registrationsucces := make(map[string]interface{})
		registrationsucces["type"] = "event"
		registrationsucces["subtype"] = "auth"
		registrationsucces["event"] = "registrationSucces"
		user := make(map[string]interface{})
		user["user"] = c.user.userJson()
		registrationsucces["data"] = user
		return registrationsucces
	}
	//TODO handle error

	return ""
}

func login(c *client, m map[string]interface{}) interface{} {
	d := m["data"]
	data := d.(map[string]interface{})
	for _, user := range c.hub.registerdUsers {
		if user.email == data["user"].(string) {
			if user.password != data["password"].(string) {
				//TODO error; wrong password
				return ""
			} else {
				c.user = *user
				authsucces := make(map[string]interface{})
				authsucces["type"] = "event"
				authsucces["subtype"] = "auth"
				authsucces["event"] = "authSucces"
				userJson := make(map[string]interface{})
				userJson["user"] = c.user.userJson()
				authsucces["user"] = userJson
				return authsucces
			}
		}
	}
	//TODO error; user not found
	return ""
}
