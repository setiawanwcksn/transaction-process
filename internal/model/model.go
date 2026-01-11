package model

import "time"

type Transaction struct {
	Timestamp    time.Time
	Counterparty string
	Type         string
	Amount       int64
	Status       string
	Description  string
}

type Statement struct {
	UploadID     string
	Balance      int64
	Transactions []Transaction
	Issues       []Transaction
	Done         bool
	Lines        int
}

type Event struct {
	UploadID string
	Tx       Transaction
	ID       string
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	NextOffset int `json:"next_offset"`
}
