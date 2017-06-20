package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (h *hub) updateDB(c client) {
	db, err := bolt.Open("chat.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Call updateDB")
	db.Update(func(tx *bolt.Tx) error {
		log.Println("Inside Bold dbj")
		b := tx.Bucket([]byte("USERS"))
		log.Println("-----------------------------------------------------------------------", c.user)
		dbuser, _ := json.Marshal(c.user)
		log.Println("--------------------------------------------------------", dbuser)
		err := b.Put([]byte(c.user.name), []byte(dbuser))
		log.Println("----------------------------------------------------------------------------", err)
		return err
	})
}

type hub struct {
	join           chan *client
	leave          chan *client
	clients        map[*client]bool
	rooms          map[string]*room
	registerdUsers map[string]*user
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
	db, err := bolt.Open("chat.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("USERS"))
		if err != nil {
			return err
		}
		return nil
	})
	users := make(map[string]*user)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("USERS"))
		b.ForEach(func(k, v []byte) error {
			var username string
			var u user
			json.Unmarshal(k, &username)
			json.Unmarshal(v, &u)
			users[username] = &u
			return nil
		})
		return nil
	})
	defer db.Close()

	return &hub{
		join:           make(chan *client),
		leave:          make(chan *client),
		clients:        make(map[*client]bool),
		rooms:          make(map[string]*room),
		registerdUsers: users,
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
			log.Println("Client joined")
		case client := <-h.leave:
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
