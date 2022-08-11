package model

type RepositoryGetOption struct {
	Offset uint64
	Limit  uint64
	Filter func(log *Log) []*Log
}
