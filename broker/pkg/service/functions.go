package service

import (
	"broker/pkg"
	"crypto/rsa"
	"errors"
	"github.com/ipaqsa/netcom/configurator"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"strings"
	"time"
)

func (broker *BrokerT) initServices() {
	for _, service := range services {
		broker.initService(service)
	}
}
func (broker *BrokerT) initService(service string) {
	broker.queues[service] = NewQueue()
}

func (broker *BrokerT) setKey() error {
	agentkey, err := broker.getAgentKey()
	if err != nil {
		return errors.New("agent`s not available")
	}

	pack := packUtils.CreatePack("broker"+":"+pkg.Config.Keyword, cryptoUtils.StringPublic(&broker.key.PublicKey))
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, broker.key, agentkey)

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetServiceKey", pack, opt)
	if err != nil {
		return err
	}

	priv := cryptoUtils.ParsePrivate(resp.Body.Data)
	if priv == nil {
		return errors.New(resp.Body.Data)
	}
	broker.key = priv
	return nil
}
func (broker *BrokerT) getAgentKey() (*rsa.PublicKey, error) {
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

func (broker *BrokerT) SelfKey() *rsa.PrivateKey {
	return broker.key
}

func (broker *BrokerT) PutReport(report *packUtils.Package) {
	infoLogger.Printf("report #%s\n", report.Head.Title)
	broker.reports.mtx.Lock()
	defer broker.reports.mtx.Unlock()
	broker.reports.data[report.Head.Title] = report
}

func (broker *BrokerT) GetReport(id string) (*packUtils.Package, error) {
	report, ok := broker.reports.data[id]
	if ok {
		return report, nil
	}
	delete(broker.reports.data, id)
	return nil, errors.New("task`s not found")
}

func (broker *BrokerT) PutMail(pack *packUtils.Package) (*packUtils.Package, error) {
	splits := strings.Split(pack.Head.Title, ".")
	if len(splits) != 2 {
		return nil, errors.New("wrong title format")
	}
	_task := splits[1]
	_service := splits[0]
	infoLogger.Printf("mail %s - service: %s\n", _task, _service)
	task := broker.createTask(_task, pack.Body.Data, _service, pack.Body.Date, pack.Head.Meta, pack.Body.Sign)
	err := broker.saveTask(task)
	if err != nil {
		return nil, err
	}
	return broker.monitor(task.Id), nil
}

func (broker *BrokerT) monitor(id string) *packUtils.Package {
	infoLogger.Printf("wait for report #%s", id)
	for {
		rep := broker.inReports(id)
		if rep != nil {
			infoLogger.Println("report`s ready")
			return rep
		}
	}
}

func (broker *BrokerT) createTask(task, data, service, date, meta, hash string) Task {
	return Task{
		Id:      SHA1(hash + data + service + date),
		From:    configurator.Info.Address,
		Service: service,
		Task:    task,
		Data:    data,
		Meta:    meta,
	}
}
func (broker *BrokerT) saveTask(task Task) error {
	queue, ok := broker.queues[task.Service]
	if ok {
		queue.Push(task)
		return nil
	}
	return errors.New("service`s not found")
}

func (broker *BrokerT) GetMail(pack *packUtils.Package) (string, error) {
	task, err := broker.getTask(pack.Head.Title)
	if err != nil {
		return "", err
	}
	marshalTask := taskMarshal(&task)
	if marshalTask == "" || task.Task == "" {
		return "", errors.New("marshal task fail")
	}
	broker.cache.data[task.Id] = task
	return marshalTask, nil
}
func (broker *BrokerT) getTask(service string) (Task, error) {
	queue, ok := broker.queues[service]
	if ok {
		if broker.queues[service].IsEmpty() {
			return Task{}, errors.New("no task")
		}
		task := queue.Pop()
		if task.Task == "" {
			return Task{}, errors.New("no task")
		}
		infoLogger.Printf("task #%s - %s | service: %s", task.Id, task.Task, task.Service)
		return task, nil
	}
	return Task{}, errors.New("service`s not found")
}

func (broker *BrokerT) cacheMonitor() {
	for {
		time.Sleep(time.Duration(pkg.Config.CacheTimeout) * time.Second)
		broker.cache.mtx.Lock()
		for key, val := range broker.cache.data {
			rep := broker.inReports(key)
			if rep == nil {
				_ = broker.saveTask(val)
			}
		}
		broker.cache.mtx.Unlock()
	}
}

func (broker *BrokerT) inReports(id string) *packUtils.Package {
	broker.reports.mtx.Lock()
	defer broker.reports.mtx.Unlock()
	for key, rep := range broker.reports.data {
		if key == id {
			delete(broker.cache.data, key)
			return rep
		}
	}
	return nil
}
