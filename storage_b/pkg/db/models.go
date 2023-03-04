package db

import (
	"database/sql"
	"sync"
)

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

type Message struct {
	Date  string `json:"date"`
	Data  string `json:"data"`
	Type  string `json:"type"`
	Meta  string `json:"meta"`
	Check string `json:"check"`
}

type Messages struct {
	Data []Message `json:"data"`
}
