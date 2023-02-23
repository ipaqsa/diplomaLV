package service

import (
	"agent/pkg"
	"agent/pkg/db"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"sync"
	"time"
)

func (agent *Agent) remove(name string) {
	agent.storage.mtx.Lock()
	defer agent.storage.mtx.Unlock()
	if _, ok := agent.storage.data[name]; ok {
		delete(agent.storage.data, name)
		err := db.RemoveKey(name)
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
	}
}
func (agent *Agent) save(name string) error {
	key := agent.Get(name)
	if key != nil {
		return errors.New("key exist")
	}
	agent.storage.mtx.Lock()
	agent.storage.data[name] = cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE)
	agent.storage.mtx.Unlock()
	return nil
}
func (agent *Agent) Get(name string) *rsa.PrivateKey {
	agent.storage.mtx.Lock()
	defer agent.storage.mtx.Unlock()
	if val, ok := agent.storage.data[name]; ok {
		return val
	} else {
		return nil
	}
}

func (agent *Agent) fetch() {
	for {
		var data []db.InsertData
		for key, val := range agent.storage.data {
			tmp := db.InsertData{
				Key:   key,
				Value: cryptoUtils.StringPrivate(val),
			}
			data = append(data, tmp)
		}
		err := db.InsertKeys(data)
		if err != nil {
			errorLogger.Printf("%s", err.Error())
			return
		}
		time.Sleep(time.Minute * 5)
	}
}

func (agent *Agent) load() error {
	data, err := db.GetKeys()
	if err != nil {
		return err
	}
	agent.storage.mtx.Lock()
	for key, val := range data {
		agent.storage.data[key] = cryptoUtils.ParsePrivate(val)
	}
	agent.storage.mtx.Unlock()
	return nil
}

func (agent *Agent) mail() {
	for {
		agent.getMail()
	}
}
func (agent *Agent) getMail() {
	pack := packUtils.CreatePack("agent", "")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, agent.Key, &agent.Get("broker").PublicKey)
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
				task.handle(agent)
			}
		}
	}
}

func (agent *Agent) sendReport(address string, pack *packUtils.Package) {
	infoLogger.Println("send report")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, agent.Get("agent"), &agent.Get("broker").PublicKey)
	pack, err := rpc.Send(address, "ServerBroker.PutReport", pack, opt)
	if err != nil {
		errorLogger.Println(err.Error())
	}
}

func (task *Task) handle(agent *Agent) {
	switch task.Task {
	case save:
		infoLogger.Printf("save #%s - %s", task.Id, task.Data)
		task.saveTask(agent)
		break
	case remove:
		infoLogger.Printf("remove #%s - %s", task.Id, task.Data)
		task.removeTask(agent)
		break
	case get:
		infoLogger.Printf("get #%s - %s", task.Id, task.Data)
		task.getTask(agent)
		break
	default:
		task.unknownTask(agent)
		break
	}
}

func (task *Task) saveTask(agent *Agent) {
	err := agent.save(task.Data)
	if err != nil {
		pack := packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		agent.sendReport(task.From, pack)
		return
	}
	pack := packUtils.CreatePack(task.Id, "ok")
	agent.sendReport(task.From, pack)
}
func (task *Task) removeTask(agent *Agent) {
	agent.remove(task.Data)
	pack := packUtils.CreatePack(task.Id, "ok")
	agent.sendReport(task.From, pack)
}
func (task *Task) getTask(agent *Agent) {
	if task.Meta != pkg.Config.KEYWORD {
		pack := packUtils.CreatePack(task.Id, "not allowed")
		pack.Head.Meta = "error"
		agent.sendReport(task.From, pack)
		return
	}
	key := agent.Get(task.Data)
	if key == nil {
		pack := packUtils.CreatePack(task.Id, "key`s not found")
		pack.Head.Meta = "error"
		agent.sendReport(task.From, pack)
		return
	}
	pack := packUtils.CreatePack(task.Id, cryptoUtils.StringPublic(&key.PublicKey))
	println(pack.Body.Data)
	agent.sendReport(task.From, pack)
}
func (task *Task) unknownTask(agent *Agent) {
	pack := packUtils.CreatePack(task.Id, "unknown task")
	agent.sendReport(task.From, pack)
}
