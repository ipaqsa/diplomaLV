package server

import (
	"encoding/json"
	"net/http"
	"node/pkg/service"
)

func parseLogin(data []byte) (*LoginForm, error) {
	var reg LoginForm
	err := json.Unmarshal(data, &reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func parseRegister(data []byte) (*RegisterForm, error) {
	var reg RegisterForm
	err := json.Unmarshal(data, &reg)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func sendAnswer(w http.ResponseWriter, data string, error string) {
	ans := Answer{
		Data:  data,
		Error: error,
	}
	jans, _ := json.Marshal(ans)
	w.Write(jans)
}

func toContactsHTML(contacts []service.User) *ContactsToHTML {
	var data ContactsToHTML
	for _, person := range contacts {
		if person.Login == service.Node.Person.Login {
			continue
		}
		data.Contact = append(data.Contact, person)
	}
	return &data
}
