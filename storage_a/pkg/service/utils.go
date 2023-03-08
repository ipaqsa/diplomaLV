package service

import (
	"encoding/json"
	"storage_a/pkg/db"
)

func marshalPerson(person *db.Person) string {
	jp, _ := json.Marshal(person)
	return string(jp)
}

func unmarshalPerson(data string) (*db.Person, error) {
	person := db.Person{}
	err := json.Unmarshal([]byte(data), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func marshalUsers(users []db.User) string {
	marshal, _ := json.Marshal(users)
	return string(marshal)
}

func unmarshalUserRequest(data string) (*UsersRequest, error) {
	var ureq UsersRequest
	err := json.Unmarshal([]byte(data), &ureq)
	if err != nil {
		return nil, err
	}
	return &ureq, nil
}
