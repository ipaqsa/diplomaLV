package service

import (
	"crypto/rsa"
	"github.com/ipaqsa/netcom/packUtils"
	"sync"
)

var services = []string{"agent", "storage_a", "storage_b", "admin"}

type BrokerT struct {
	cache   Cache
	key     *rsa.PrivateKey
	queues  map[string]*Queue
	reports Reports
}

type Cache struct {
	data map[string]Task
	mtx  sync.Mutex
}

type Queue struct {
	data []Task
	mtx  sync.Mutex
}

type Reports struct {
	data map[string]*packUtils.Package
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
