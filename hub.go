package main

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type hub struct {
	join           chan *client
	leave          chan *client
	clients        map[*client]bool
	rooms          map[string]*room
	registerdUsers map[string]*user
	db             *mgo.Session
}

type dbuser struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"password" bson:"password"`
	AvatarURL string        `json:"avatarURL" bson:"avatarURL"`
}

func (h *hub) createUserEmailMap() map[string]bool {
	emails := make(map[string]bool)
	for _, user := range h.registerdUsers {
		emails[user.email] = true
	}
	return emails
}

func (u *dbuser) castUser() *user {
	return &user{
		id:        u.Id,
		name:      u.Name,
		email:     u.Email,
		password:  u.Password,
		avatarURL: u.AvatarURL,
	}
}

func (h *hub) sendToAll(msg interface{}) {
	for k := range h.clients {
		k.send <- msg
	}
}

func (h hub) getAllRegisteredUserNames() []string {
	userNames := make([]string, 0, 0)
	for _, v := range h.registerdUsers {
		userNames = append(userNames, v.name)
	}
	return userNames
}
func (h hub) usersJson() []map[string]interface{} {
	users := make([]map[string]interface{}, 0, 0)
	for k := range h.clients {
		users = append(users, k.user.userJson())
	}
	return users
}

func (h hub) roomsJson() []map[string]interface{} {
	rooms := make([]map[string]interface{}, 0, 0)
	for _, k := range h.rooms {
		rooms = append(rooms, k.roomJson())
	}
	return rooms
}

func newHub() *hub {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	c := session.DB("chat").C("users")
	users := make(map[string]*user)
	var mgouser dbuser
	iter := c.Find(nil).Iter()
	for iter.Next(&mgouser) {
		users[mgouser.Name] = mgouser.castUser()
	}
	iter.Close()

	return &hub{
		join:           make(chan *client),
		leave:          make(chan *client),
		clients:        make(map[*client]bool),
		rooms:          make(map[string]*room),
		registerdUsers: users,
		db:             session,
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
			log.Println("Client joined")
		case client := <-h.leave:
			msg := userLeaves(client)
			h.sendToAll(msg)
			delete(h.clients, client)
			close(client.send)
			log.Println("Client left")
		}
	}
}

const (
	socketBuffersize  = 1024
	messageBuffersize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBuffersize,
	WriteBufferSize: socketBuffersize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *hub) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("ServeHTTP:", err)
		return
	}
	client := &client{
		hub:    h,
		socket: socket,
		send:   make(chan interface{}),
	}
	h.join <- client
	defer func() { h.leave <- client }()
	go client.write()
	client.read()
}
