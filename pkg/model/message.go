package model

type ExportLogsResponse struct {
	Logs      []*Log    `json:"logs"`
	NextToken NextToken `json:"next_token"`
}

type RepositoryGetInput struct {
	Limit  uint64
	Token  NextToken
	Filter func(log *Log) []*Log
}

type RepositoryGetOutput struct {
	Logs      []*Log
	NextToken NextToken
}
