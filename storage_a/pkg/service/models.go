package service

import (
	"crypto/rsa"
	"storage_a/pkg/db"
)

const (
	remove = "remove"
	save   = "save"
	auth   = "auth"
	users  = "users"
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