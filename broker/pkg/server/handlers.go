package server

import (
	"broker/pkg"
	"broker/pkg/service"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
)

func (server *ServerBroker) GetMail(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, service.Broker.SelfKey(), nil)
	pack, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	mail, err := service.Broker.GetMail(pack)
	if err != nil {
		if err.Error() != "no task" {
			errorLogger.Println(err.Error())
		}
		return err
	}
	ans.Body.Data = mail
	rpc.SendAnswer(ans, opt, response)
	return nil
}

func (server *ServerBroker) PutMail(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, service.Broker.SelfKey(), nil)
	pack, _, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	pack, err = service.Broker.PutMail(pack)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	rpc.SendAnswer(pack, opt, response)
	return nil
}

func (server *ServerBroker) GetReport(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, service.Broker.SelfKey(), nil)
	pack, _, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	report, err := service.Broker.GetReport(pack.Head.Title)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	rpc.SendAnswer(report, opt, response)
	return nil
}

func (server *ServerBroker) PutReport(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, service.Broker.SelfKey(), nil)
	pack, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	service.Broker.PutReport(pack)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	rpc.SendAnswer(ans, opt, response)
	return nil
}

func (server *ServerBroker) Health(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(false, 0, nil, nil)
	_, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	ans.Body.Data = "alive"
	rpc.SendAnswer(ans, opt, response)
	return nil
}
