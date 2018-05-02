package main

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clinets
	forward chan []byte
}
