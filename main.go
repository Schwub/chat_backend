package main

import (
	"log"
	"net/http"
)

func main() {
	h := newHub()
	log.Println("Created Hub")
	http.Handle("/", h)
	log.Println("Handle HTTP Request")
	go h.run()
	defer h.db.Close()
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Println("ListenAndServe:", err)
	}
}
