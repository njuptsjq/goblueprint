package main

import (
	"goblueprint/trace"
	"log"
	"net/http"

	"github.com/stretchr/objx"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clinets
	forward chan *message
	// join  is a channel for clients wishing to join
	join chan *client
	// leave is a channel for leaving
	leave chan *client
	// clients holds all current clients in this room
	clients map[*client]bool

	tracer trace.Tracer

	avatar Avatar
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			// forward message
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- sent to client", msg.Message)
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
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}
