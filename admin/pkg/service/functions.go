package service

import (
	"admin/pkg"
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

func (admin *AdminT) setKey() error {
	agentkey, err := admin.getAgentKey()
	if err != nil {
		return errors.New("agent`s not available")
	}

	pack := packUtils.CreatePack("broker"+":"+pkg.Config.Keyword, cryptoUtils.StringPublic(&admin.key.PublicKey))
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, admin.key, agentkey)

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetServiceKey", pack, opt)
	if err != nil {
		return err
	}

	priv := cryptoUtils.ParsePrivate(resp.Body.Data)
	if priv == nil {
		return errors.New(resp.Body.Data)
	}
	admin.key = priv
	return nil
}
func (admin *AdminT) getAgentKey() (*rsa.PublicKey, error) {
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
func (admin *AdminT) getBrokerKey() (*rsa.PublicKey, error) {
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

func (admin *AdminT) SelfKey() *rsa.PrivateKey {
	return admin.key
}

func (admin *AdminT) mail() {
	for {
		admin.getMail()
	}
}

func (admin *AdminT) getMail() {
	brokerKey, err := admin.getBrokerKey()
	if err != nil {
		errorLogger.Println(err.Error())
		return
	}
	pack := packUtils.CreatePack("admin", "")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, admin.key, brokerKey)
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
				task.handle(admin, brokerKey)
			}
		}
	}
}

func (admin *AdminT) sendReport(address string, pack *packUtils.Package, brokerKey *rsa.PublicKey) {
	infoLogger.Printf("send report")
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, admin.key, brokerKey)
	pack, err := rpc.Send(address, "ServerBroker.PutReport", pack, opt)
	if err != nil {
		errorLogger.Println(err.Error())
	}
}

func (task *Task) handle(admin *AdminT, brokerKey *rsa.PublicKey) {
	switch task.Task {
	case register:
		task.registerTask(admin, brokerKey)
		break
	case auth:
		task.authTask(admin, brokerKey)
		break
	default:
		task.unknownTask(admin, brokerKey)
		break
	}
}

func (task *Task) registerTask(admin *AdminT, brokerKey *rsa.PublicKey) {
	err := admin.register(task)
	if err != nil {
		pack := packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		admin.sendReport(task.From, pack, brokerKey)
		return
	}
	pack := packUtils.CreatePack(task.Id, "ok")
	admin.sendReport(task.From, pack, brokerKey)
}
func (task *Task) authTask(admin *AdminT, brokerKey *rsa.PublicKey) {
	pack, err := admin.auth(task)
	if err != nil {
		pack = packUtils.CreatePack(task.Id, err.Error())
		pack.Head.Meta = "error"
		admin.sendReport(task.From, pack, brokerKey)
		return
	}
	admin.sendReport(task.From, pack, brokerKey)
}
func (task *Task) unknownTask(admin *AdminT, brokerKey *rsa.PublicKey) {
	pack := packUtils.CreatePack(task.Id, "unknown task")
	admin.sendReport(task.From, pack, brokerKey)
}

func (admin *AdminT) register(task *Task) error {
	infoLogger.Printf("register #%s - %s\n", task.Id, task.Meta)
	data, meta, err := admin.sendMail("storage_a", "save", task.Data, "")
	if err != nil {
		return err
	}
	if meta == "error" {
		return errors.New(data)
	}
	infoLogger.Printf("register #%s(key`s creating)- %s\n", task.Id, task.Meta)
	data, meta, err = admin.sendMail("agent", "save", task.Meta, "")
	if err != nil {
		infoLogger.Printf("rollback #%s - %s\n", task.Id, task.Meta)
		_, _, err = admin.sendMail("storage_a", "remove", task.Meta, "")
		if err != nil {
			errorLogger.Println(err.Error())
			return err
		}
		return err
	}
	if meta == "error" {
		infoLogger.Printf("rollback #%s - %s\n", task.Id, task.Meta)
		_, _, err = admin.sendMail("storage_a", "remove", task.Meta, "")
		if err != nil {
			errorLogger.Println(err.Error())
			return err
		}
		return errors.New(data)
	}
	return nil
}
func (admin *AdminT) auth(task *Task) (*packUtils.Package, error) {
	infoLogger.Printf("auth #%s - %s\n", task.Id, task.Data)
	data, meta, err := admin.sendMail("storage_a", "auth", task.Data, task.Meta)
	if err != nil {
		errorLogger.Println(err.Error())
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	pack := packUtils.CreatePack(task.Id, data)
	infoLogger.Printf("auth #%s(key`s getting) - %s\n", task.Id, task.Data)
	data, meta, err = admin.sendMail("agent", "get", task.Data, pkg.Config.Keyword)
	if err != nil {
		errorLogger.Println(err.Error())
		return nil, err
	}
	if meta == "error" {
		return nil, errors.New(data)
	}
	pack.Head.Meta = data
	return pack, nil
}

func (admin *AdminT) sendMail(service, task, data, meta string) (string, string, error) {
	brokerKey, err := admin.getBrokerKey()
	if err != nil {
		return "", "", err
	}
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, admin.SelfKey(), brokerKey)
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
