package service

import (
	"crypto/rsa"
	"storage_b/pkg/db"
)

const (
	remove = "remove"
	save   = "save"
	get    = "get"
	count  = "count"
)

type StorageT struct {
	db  *db.DB
	key *rsa.PrivateKey
}

type Task struct {
	From    string `json:"from"`
	Id      string `json:"id"`
	Service string `json:"service"`
	Task    string `json:"task"`
	Data    string `json:"data"`
	Meta    string `json:"meta"`
}
