package service

import "crypto/rsa"

type NodeT struct {
	Person *Person
	Key    *rsa.PrivateKey
	Status bool
}

type Person struct {
	Login     string `json:"login"`
	Hash      string `json:"hash"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Room      int    `json:"room"`
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

type User struct {
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type FileMessage struct {
	Title string `json:"title"`
	Data  string `json:"data"`
	Meta  string `json:"meta"`
}
