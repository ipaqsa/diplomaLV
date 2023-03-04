package service

import (
	"crypto/rsa"
)

const (
	remove = "remove"
	save   = "save"
	get    = "get"
)

type StorageT struct {
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

type FileMessage struct {
	Title string `json:"title"`
	Data  string `json:"data"`
	Meta  string `json:"meta"`
}
