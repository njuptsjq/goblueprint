package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clinets
	forward chan []byte
	// join  is a channel for clients wishing to join
	join chan *client
	// leave is a channel for leaving
	leave chan *client
	// clinets holds all current clients in this room
	clinets map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clinets[client] = true
		case client := <-r.leave:
			// leaving
			delete(r.clinets, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message
			for client := range r.clinets {
				client.send <- msg
			}

		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
