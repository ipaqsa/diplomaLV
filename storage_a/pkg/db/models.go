package db

import (
	"database/sql"
	"sync"
)

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

type Person struct {
	Login     string `json:"login"`
	Hash      string `json:"hash"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Room      int    `json:"room"`
}

type User struct {
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Count     int    `json:"count"`
}
