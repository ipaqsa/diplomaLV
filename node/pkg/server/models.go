package server

import "node/pkg/service"

type RegisterForm struct {
	Login      string `json:"login"`
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	Password   string `json:"password"`
	Room       string `json:"room"`
}

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Answer struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type DataToHTML struct {
	Receiver string           `json:"Receiver"`
	Contacts *ContactsToHTML  `json:"Contacts"`
	Messages service.Messages `json:"Messages"`
}

type ContactToHTML struct {
	Login      string
	FirstName  string
	SecondName string
}

type ContactsToHTML struct {
	Contact []ContactToHTML
}

type sendFromHTML struct {
	Data string `json:"data"`
}
