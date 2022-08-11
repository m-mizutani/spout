package model

import "time"

type Message struct {
	Log   *Log
	Error error
}

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Tag       string    `json:"tag"`
	Data      any       `json:"data"`
}
