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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("ListenAndServe:", err)
	}
}
