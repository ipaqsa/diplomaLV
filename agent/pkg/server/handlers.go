package server

import (
	"agent/pkg"
	"agent/pkg/service"
	"errors"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"strings"
)

func (server *ServerAgent) GetSelfKey(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(false, 0, nil, nil)
	_, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	ans.Body.Data = cryptoUtils.StringPublic(&service.Service.Key.PublicKey)
	rpc.SendAnswer(ans, opt, response)
	return nil
}

func (server *ServerAgent) GetServiceKey(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, service.Service.Key, nil)
	pack, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	splits := strings.Split(pack.Head.Title, ":")
	if len(splits) != 2 {
		errorLogger.Println("wrong title")
		return errors.New("wrong title")
	}
	if splits[1] != pkg.Config.KEYWORD {
		errorLogger.Println("wrong keyword")
		return errors.New("wrong keyword")
	}

	key := service.Service.Get(splits[0])
	if key == nil {
		errorLogger.Println("key`s nil")
		return errors.New("key`s nil")
	}
	ans.Body.Data = cryptoUtils.StringPrivate(key)
	rpc.SendAnswer(ans, opt, response)
	return nil
}

func (server *ServerAgent) GetBrokerKey(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(false, 0, nil, nil)
	_, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	ans.Body.Data = cryptoUtils.StringPublic(&service.Service.Get("broker").PublicKey)
	rpc.SendAnswer(ans, opt, response)
	return nil
}
