package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/ipaqsa/netcom/configurator"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"node/pkg"
	"node/pkg/service"
	"os"
	"strconv"
)

func init() {
	err := configurator.InitConfiguration(&pkg.Config, "0.0.1")
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	configurator.InitInfo(pkg.Config.Port)
}

func main() {
	service.NewNode()
	//err := Register("stefan", "stefan", "Stepan", "Nashville", 2)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//err = Register("alice", "alice", "Alice", "Wonderland", 2)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	_, key, err := Auth("stefan", "stefan")
	if err != nil {
		println(err.Error())
		return
	}
	service.Node.Key = key
	err = Send("Hello Alice", &service.Node.Key.PublicKey)
	if err != nil {
		println(err.Error())
		return
	}
	//usrs, err := Users(2)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//for _, usr := range usrs {
	//	println(usr.Login)
	//}
	//auth, err := Auth("rtefan", "Stefan")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//println(auth.Lastname)
}

func Register(login, password, firstname, lastname string, room int) error {
	person := service.Person{Login: login, Hash: cryptoUtils.Base64Encode(cryptoUtils.HashSum([]byte(password))), Firstname: firstname, Lastname: lastname, Room: room}
	jp := service.MarshalPerson(&person)
	data, meta, err := service.Node.SendMail("admin", "register", jp, person.Login)
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}

func Auth(login, password string) (*service.Person, *rsa.PrivateKey, error) {
	data, meta, err := service.Node.SendMail("admin", "auth", login, cryptoUtils.Base64Encode(cryptoUtils.HashSum([]byte(password))))
	if err != nil {
		return nil, nil, err
	}
	if meta == "error" {
		return nil, nil, errors.New(data)
	}
	person, err := service.UnmarshalPerson(data)
	if err != nil {
		return nil, nil, err
	}
	println(meta)
	key := cryptoUtils.ParsePrivate(meta)
	if key == nil {
		return nil, nil, errors.New("key parse error")
	}
	return person, key, nil
}

func Users(room int) ([]service.User, error) {
	data, meta, err := service.Node.SendMail("storage_a", "users", strconv.Itoa(room), pkg.Config.Keyword)
	if err != nil {
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	users := service.UnmarshalUsers(data)
	if users == nil {
		return nil, errors.New("unmarshal fail")
	}
	return users, nil
}

func Send(data string, key *rsa.PublicKey) error {
	msg := service.Message{Data: data, Meta: cryptoUtils.StringPublic(key)}
	jmsg, _ := json.Marshal(msg)
	data, meta, err := service.Node.SendMail("storage_b", "save", string(jmsg), "ALICEHASH")
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	msg = service.Message{Data: data, Meta: cryptoUtils.StringPublic(key)}
	jmsg, _ = json.Marshal(msg)
	data, meta, err = service.Node.SendMail("storage_b", "save", string(jmsg), cryptoUtils.StringPublic(key))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
