package service

import (
	"crypto/rsa"
	"sync"
)

const (
	remove  = "remove"
	save    = "save"
	get     = "get"
	getPriv = "getpriv"
)

type Agent struct {
	Key     *rsa.PrivateKey
	storage Storage
}

type Storage struct {
	data map[string]*rsa.PrivateKey
	mtx  sync.Mutex
}

type Task struct {
	From    string `json:"from"`
	Id      string `json:"id"`
	Service string `json:"service"`
	Task    string `json:"task"`
	Data    string `json:"data"`
	Meta    string `json:"meta"`
}
