package service

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"node/pkg"
	"strconv"
	"time"
)

func (node *NodeT) getAgentKey() (*rsa.PublicKey, error) {
	pack := packUtils.CreatePack("GET", "")

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetSelfKey", pack, nil)
	if err != nil {
		return nil, err
	}

	agentkey := cryptoUtils.ParsePublic(resp.Body.Data)
	if agentkey == nil {
		return nil, errors.New("agent key parse")
	}
	return agentkey, nil
}
func (node *NodeT) getBrokerKey() (*rsa.PublicKey, error) {
	pack := packUtils.CreatePack("GET", "")

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetBrokerKey", pack, nil)
	if err != nil {
		return nil, err
	}

	agentkey := cryptoUtils.ParsePublic(resp.Body.Data)
	if agentkey == nil {
		return nil, errors.New("broker key parse")
	}
	return agentkey, nil
}

func (node *NodeT) SelfKey() *rsa.PrivateKey {
	return node.Key
}

func (node *NodeT) SendMail(service, task, data, meta string) (string, string, error) {
	broker, err := node.getBrokerKey()
	if err != nil {
		return "", "", err
	}
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, node.SelfKey(), broker)
	pack := packUtils.CreatePack(service+"."+task, data)
	pack.Head.Meta = meta

	resp, err := rpc.Send(pkg.Config.Ingress, "IngressT.Broker", pack, opt)
	if err != nil {
		return "", "", err
	}
	return resp.Body.Data, resp.Head.Meta, nil
}

func (node *NodeT) Register(login, password, firstname, lastname string, room int) error {
	person := Person{Login: login, Hash: cryptoUtils.Base64Encode(cryptoUtils.HashSum([]byte(password))), Firstname: firstname, Lastname: lastname, Room: room}
	jp := MarshalPerson(&person)
	data, meta, err := Node.SendMail("admin", "register", jp, person.Login)
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
func (node *NodeT) Authentication(login, password string) error {
	data, meta, err := node.SendMail("admin", "auth", login, cryptoUtils.Base64Encode(cryptoUtils.HashSum([]byte(password))))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	person, err := UnmarshalPerson(data)
	if err != nil {
		return err
	}
	key := cryptoUtils.ParsePrivate(meta)
	if key == nil {
		return errors.New("key parse error")
	}
	node.Person = person
	node.Key = key
	node.Switch()
	return nil
}
func (node *NodeT) RemoveAccount() error {
	data, meta, err := node.SendMail("admin", "remove", node.Person.Login, pkg.Config.Keyword)
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
func (node *NodeT) GetContacts() ([]User, error) {
	ureq := UsersRequest{SenderKey: cryptoUtils.StringPublic(&node.Key.PublicKey), Room: strconv.Itoa(node.Person.Room)}
	jureq, err := MarshalUsersRequest(&ureq)
	if err != nil {
		return nil, err
	}
	data, meta, err := node.SendMail("storage_a", "users", jureq, pkg.Config.Keyword)
	if err != nil {
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	users := UnmarshalUsers(data)
	if users == nil {
		return nil, errors.New("unmarshal fail")
	}
	return users, nil
}
func (node *NodeT) Send(text, receiver string) error {
	receiverKey, err := node.getKey(receiver)
	if err != nil {
		return err
	}
	if receiverKey == nil {
		return errors.New("key parse error")
	}
	date := time.Now().Format(time.RFC822)
	msg := Message{Data: text, Meta: cryptoUtils.StringPublic(&node.Key.PublicKey), Date: date, Type: "text"}
	jmsg, _ := json.Marshal(msg)
	data, meta, err := Node.SendMail("storage_b", "save", string(jmsg), cryptoUtils.StringPublic(receiverKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	data, meta, err = Node.SendMail("storage_b", "save", string(jmsg), cryptoUtils.StringPublic(receiverKey)+","+cryptoUtils.StringPublic(&node.Key.PublicKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
func (node *NodeT) SendFile(file []byte, filename, receiver string) error {
	jf, err := MarshalFile(file, filename, cryptoUtils.StringPublic(&node.Key.PublicKey))
	if err != nil {
		return err
	}
	receiverKey, err := node.getKey(receiver)
	if err != nil {
		return err
	}
	if receiverKey == nil {
		return errors.New("key parse error")
	}
	data, meta, err := node.SendMail("storage_f", "save", jf, cryptoUtils.StringPublic(receiverKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	data, meta, err = node.SendMail("storage_f", "save", jf, cryptoUtils.StringPublic(receiverKey)+","+cryptoUtils.StringPublic(&node.Key.PublicKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	date := time.Now().Format(time.RFC822)
	msg := Message{Data: filename, Meta: cryptoUtils.StringPublic(&node.Key.PublicKey), Date: date, Type: "file"}
	jmsg, _ := json.Marshal(msg)
	data, meta, err = Node.SendMail("storage_b", "save", string(jmsg), cryptoUtils.StringPublic(receiverKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	data, meta, err = Node.SendMail("storage_b", "save", string(jmsg), cryptoUtils.StringPublic(receiverKey)+","+cryptoUtils.StringPublic(&node.Key.PublicKey))
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
func (node *NodeT) GetFile(filename, receiver string) (*FileMessage, error) {
	println(filename)
	file := &FileMessage{filename, "", cryptoUtils.StringPublic(&node.Key.PublicKey)}
	jf, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}
	receiverKey, err := node.getKey(receiver)
	if err != nil {
		return nil, err
	}
	if receiverKey == nil {
		return nil, errors.New("key parse error")
	}
	data, meta, err := node.SendMail("storage_f", "get", string(jf), cryptoUtils.StringPublic(receiverKey))
	if err != nil {
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	file, err = UnmarshalFile(data)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (node *NodeT) Messages(receiver string) (*Messages, error) {
	receiverKey, err := node.getKey(receiver)
	if err != nil {
		return nil, err
	}
	if receiverKey == nil {
		return nil, errors.New("key parse error")
	}
	data, meta, err := Node.SendMail("storage_b", "get", cryptoUtils.StringPublic(&node.Key.PublicKey), cryptoUtils.StringPublic(receiverKey))
	if err != nil {
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	return UnmarshalMessages(data), nil
}
func (node *NodeT) Update(passwordChange string) error {
	person := Person{Login: node.Person.Login, Hash: node.Person.Hash, Firstname: node.Person.Firstname, Lastname: node.Person.Lastname, Room: node.Person.Room}
	jp := MarshalPerson(&person)
	data, meta, err := Node.SendMail("admin", "update", jp, passwordChange)
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	return nil
}
func (node *NodeT) getKey(login string) (*rsa.PublicKey, error) {
	data, meta, err := node.SendMail("agent", "get", login, pkg.Config.Keyword)
	if err != nil {
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	return cryptoUtils.ParsePublic(data), nil
}

func (node *NodeT) Switch() {
	if node.Status {
		infoLogger.Println("status: no auth")
		node.Key = cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE)
		node.Person = &Person{}
		node.Status = false
	} else {
		infoLogger.Println("status: auth")
		node.Status = true
	}
}
