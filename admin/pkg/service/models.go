package service

import (
	"crypto/rsa"
)

const (
	register = "register"
	auth     = "auth"
)

type AdminT struct {
	name string
	key  *rsa.PrivateKey
}

type Task struct {
	From    string `json:"from"`
	Id      string `json:"id"`
	Service string `json:"service"`
	Task    string `json:"task"`
	Data    string `json:"data"`
	Meta    string `json:"meta"`
}
