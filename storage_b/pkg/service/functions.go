package service

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"storage_b/pkg"
	"storage_b/pkg/db"
	"strings"
	"sync"
	"time"
)

func (store *StorageT) setKey() error {
	agentkey, err := store.getAgentKey()
	if err != nil {
		return errors.New("agent`s not available")
	}

	pack := packUtils.CreatePack("broker"+":"+pkg.Config.Keyword, cryptoUtils.StringPublic(&store.key.PublicKey))
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, store.key, agentkey)

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetServiceKey", pack, opt)
	if err != nil {
		return err
	}

	priv := cryptoUtils.ParsePrivate(resp.Body.Data)
	if priv == nil {
		return errors.New(resp.Body.Data)
	}
	store.key = priv
	return nil
}
func (store *StorageT) getAgentKey() (*rsa.PublicKey, error) {
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
func (store *StorageT) getBrokerKey() (*rsa.PublicKey, error) {
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

func (store *StorageT) SelfKey() *rsa.PrivateKey {
	return store.key
}

func (store *StorageT) getKey(login string) *rsa.PublicKey {
	data, _, err := store.sendMail("agent", "get", login, "")
	if err != nil {
		errorLogger.Println(err.Error())
		return nil
	}
	return cryptoUtils.ParsePublic(data)
}

func (store *StorageT) sendMail(service, task, data, meta string) (string, string, error) {
	broker, err := store.getBrokerKey()
	if err != nil {
		return "", "", err
	}
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, store.SelfKey(), broker)
	pack := packUtils.CreatePack(service+"."+task, data)
	pack.Head.Meta = meta
	for _, broker := range pkg.Config.Brokers {
		resp, err := rpc.Send(broker, "ServerBroker.PutMail", pack, opt)
		if err != nil {
			errorLogger.Println(err.Error())
			continue
		}
		return resp.Body.Data, resp.Head.Meta, nil
	}
	return "", "", errors.New("send fail")
}

func (store *StorageT) save(task *Task) error {
	msg := db.UnmarshalMessage(task.Data)
	if msg == nil {
		return errors.New("unmarshal error")
	}
	type_save := 0

	receiver := task.Meta
	sender := msg.Meta
	if receiver == sender {
		type_save = 1
	}
	sender = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(sender))), "/", "S")
	receiver = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(receiver))), "/", "S")
	//
	sender = strings.ReplaceAll(sender, "+", "S")
	receiver = strings.ReplaceAll(receiver, "+", "S")
	//
	sender = strings.ReplaceAll(sender, "=", "S")
	receiver = strings.ReplaceAll(receiver, "=", "S")
	//
	sender = strings.ReplaceAll(sender, "?", "S")
	receiver = strings.ReplaceAll(receiver, "?", "S")
	//
	err := store.db.AddMessage(sender, receiver, msg.Data, msg.Date, type_save)
	if err != nil {
		return err
	}
	return nil
}

func (store *StorageT) get(task *Task) (*packUtils.Package, error) {
	receiver := task.Meta
	sender := task.Data

	sender = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(sender))), "/", "S")
	receiver = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(receiver))), "/", "S")

	sender = strings.ReplaceAll(sender, "+", "S")
	receiver = strings.ReplaceAll(receiver, "+", "S")

	sender = strings.ReplaceAll(sender, "=", "S")
	receiver = strings.ReplaceAll(receiver, "=", "S")

	sender = strings.ReplaceAll(sender, "?", "S")
	receiver = strings.ReplaceAll(receiver, "?", "S")

	msg := store.db.GetMessages(sender, receiver)
	if msg == nil {
		return nil, errors.New("no msg")
	}
	return packUtils.CreatePack("answer", db.MarshalMessages(msg)), nil
}

func (store *StorageT) mail() {
	for {
		store.getMail()
	}
}

func (store *StorageT) getMail() {
	brokerKey, err := store.getBrokerKey()
	if err != nil {
		errorLogger.Println(err.Error())
		return
	}
	pack := packUtils.CreatePack("storage_b", "")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, store.key, brokerKey)
	var mtx sync.Mutex
	mtx.Lock()
	defer mtx.Unlock()
	for _, addr := range pkg.Config.Brokers {
		if addr != "" {
			ans, err := rpc.Send(addr, "ServerBroker.GetMail", pack, opt)
			if err != nil {
				if err.Error() == fmt.Sprintf("dial tcp %s: connect: connection refused", addr) {
					errorLogger.Println("fail connect to broker")
					time.Sleep(time.Second)
				} else if err.Error() == "no task" {
					continue
				} else {
					errorLogger.Println(err.Error())
				}
				continue
			}
			var task Task
			err = json.Unmarshal([]byte(ans.Body.Data), &task)
			if err != nil {
				errorLogger.Println(err.Error())
				continue
			}
			if task.Task != "" {
				task.handle(store, brokerKey)
			}
		}
	}
}

func (store *StorageT) sendReport(address string, pack *packUtils.Package, brokerKey *rsa.PublicKey) {
	infoLogger.Println("send report")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, store.key, brokerKey)
	pack, err := rpc.Send(address, "ServerBroker.PutReport", pack, opt)
	if err != nil {
		errorLogger.Println(err.Error())
	}
}

func (task *Task) handle(store *StorageT, brokerKey *rsa.PublicKey) {
	switch task.Task {
	case save:
		infoLogger.Printf("save #%s", task.Id)
		task.saveTask(store, brokerKey)
		break
	//case remove:
	//	task.removeTask(store, brokerKey)
	//	break
	case get:
		infoLogger.Printf("get #%s", task.Id)
		task.getTask(store, brokerKey)
		break
	default:
		task.unknownTask(store, brokerKey)
		break
	}
}

func (task *Task) saveTask(store *StorageT, brokerKey *rsa.PublicKey) {
	err := store.save(task)
	if err != nil {
		pack := packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		store.sendReport(task.From, pack, brokerKey)
		return
	}
	pack := packUtils.CreatePack(task.Id, "ok")
	store.sendReport(task.From, pack, brokerKey)
}

func (task *Task) getTask(store *StorageT, brokerKey *rsa.PublicKey) {
	pack, err := store.get(task)
	if err != nil {
		pack = packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		store.sendReport(task.From, pack, brokerKey)
		return
	}
	store.sendReport(task.From, pack, brokerKey)
	store.sendReport(task.From, pack, brokerKey)
}

func (task *Task) unknownTask(store *StorageT, brokerKey *rsa.PublicKey) {
	pack := packUtils.CreatePack(task.Id, "unknown task")
	store.sendReport(task.From, pack, brokerKey)
}
