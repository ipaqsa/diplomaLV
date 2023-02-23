package service

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"storage_a/pkg"
	"strconv"
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

func (store *StorageT) save(task *Task) error {
	person, err := unmarshalPerson(task.Data)
	if err != nil {
		return err
	}
	err = store.db.RegisterPerson(person)
	if err != nil {
		return err
	}
	return nil
}
func (store *StorageT) auth(task *Task) (*packUtils.Package, error) {
	person, err := store.db.GetPerson(task.Data, task.Meta)
	if err != nil {
		return nil, err
	}
	sperson := marshalPerson(person)
	if sperson == "" {
		return nil, err
	}
	return packUtils.CreatePack(task.Id, sperson), nil
}
func (store *StorageT) remove(task *Task) {
	err := store.db.RemovePerson(task.Data)
	if err != nil {
		errorLogger.Println(err.Error())
		return
	}
}
func (store *StorageT) users(task *Task) (*packUtils.Package, error) {
	if task.Meta != pkg.Config.Keyword {
		return nil, errors.New("not allowed")
	}
	room, err := strconv.Atoi(task.Data)
	if err != nil {
		return nil, err
	}
	usrs := store.db.GetAllMember(room)
	if usrs == nil || len(usrs) == 0 {
		return nil, errors.New("no users")
	}
	jusrs := marshalUsers(usrs)
	if jusrs == "" {
		return nil, errors.New("marshal fail")
	}
	return packUtils.CreatePack(task.Id, jusrs), nil
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
	pack := packUtils.CreatePack("storage_a", "")
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
	infoLogger.Printf("send report")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, store.key, brokerKey)
	pack, err := rpc.Send(address, "ServerBroker.PutReport", pack, opt)
	if err != nil {
		errorLogger.Println(err.Error())
	}
}

func (task *Task) handle(store *StorageT, brokerKey *rsa.PublicKey) {
	switch task.Task {
	case users:
		infoLogger.Printf("users #%s", task.Id)
		task.usersTask(store, brokerKey)
		break
	case save:
		infoLogger.Printf("save #%s", task.Id)
		task.saveTask(store, brokerKey)
		break
	case remove:
		infoLogger.Printf("remove #%s - %s", task.Id, task.Data)
		task.removeTask(store, brokerKey)
		break
	case auth:
		infoLogger.Printf("auth #%s - %s", task.Id, task.Data)
		task.authTask(store, brokerKey)
		break
	default:
		task.unknownTask(store, brokerKey)
		break
	}
}

func (task *Task) usersTask(store *StorageT, brokerKey *rsa.PublicKey) {
	pack, err := store.users(task)
	if err != nil {
		pack = packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		store.sendReport(task.From, pack, brokerKey)
		return
	}
	store.sendReport(task.From, pack, brokerKey)
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
func (task *Task) removeTask(store *StorageT, brokerKey *rsa.PublicKey) {
	store.remove(task)
	pack := packUtils.CreatePack(task.Id, "ok")
	store.sendReport(task.From, pack, brokerKey)
}
func (task *Task) authTask(store *StorageT, brokerKey *rsa.PublicKey) {
	pack, err := store.auth(task)
	if err != nil {
		pack = packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		store.sendReport(task.From, pack, brokerKey)
		return
	}
	store.sendReport(task.From, pack, brokerKey)
}
func (task *Task) unknownTask(store *StorageT, brokerKey *rsa.PublicKey) {
	pack := packUtils.CreatePack(task.Id, "unknown task")
	store.sendReport(task.From, pack, brokerKey)
}
