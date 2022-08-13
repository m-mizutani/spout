package model

import "time"

type Message struct {
	Log   *Log
	Error error
}

type Log struct {
	ID        LogID
	Timestamp *time.Time `json:"timestamp"`
	Tag       string     `json:"tag"`
	Data      any        `json:"data"`
}

func NewLog(ts *time.Time, tag string, data any) *Log {
	return &Log{
		ID:        NewLogID(),
		Timestamp: ts,
		Tag:       tag,
		Data:      data,
	}
}
