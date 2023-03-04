package service

import (
	"encoding/json"
	"github.com/ipaqsa/netcom/cryptoUtils"
)

func MarshalPerson(person *Person) string {
	jp, _ := json.Marshal(person)
	return string(jp)
}
func UnmarshalPerson(data string) (*Person, error) {
	person := Person{}
	err := json.Unmarshal([]byte(data), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func MarshalMessages(msg *Messages) string {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(jmsg)
}
func UnmarshalMessages(data string) *Messages {
	var msg Messages
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		return nil
	}
	return &msg
}

func UnmarshalUsers(data string) []User {
	var users []User
	err := json.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil
	}
	return users
}

func MarshalFile(file []byte, filename, sender string) (string, error) {
	data := cryptoUtils.Base64Encode(file)
	fm := FileMessage{Title: filename, Data: data, Meta: sender}
	jfm, err := json.Marshal(fm)
	return string(jfm), err
}
