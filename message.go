package main

import (
	"time"
)

// message reprents a single message
type message struct {
	Name    string
	Message string
	When    time.Time
}
