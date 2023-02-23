package service

import (
	"crypto/rsa"
	"errors"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"node/pkg"
)

func (node *NodeT) setKey() error {
	agentkey, err := node.getAgentKey()
	if err != nil {
		return errors.New("agent`s not available")
	}

	pack := packUtils.CreatePack("broker"+":"+pkg.Config.Keyword, cryptoUtils.StringPublic(&node.Key.PublicKey))
	opt := rpc.CreateOptions(true, pkg.Config.SKEY_SIZE, node.Key, agentkey)

	resp, err := rpc.Send(pkg.Config.Agent, "ServerAgent.GetServiceKey", pack, opt)
	if err != nil {
		return err
	}

	priv := cryptoUtils.ParsePrivate(resp.Body.Data)
	if priv == nil {
		return errors.New(resp.Body.Data)
	}
	node.Key = priv
	return nil
}

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
