package server

import (
	"github.com/ipaqsa/netcom/packUtils"
	"github.com/ipaqsa/netcom/rpc"
	"ingress/pkg/service"
)

func (ingress *IngressT) Broker(data []byte, response *packUtils.Package) error {
	opt := rpc.CreateOptions(false, 0, nil, nil)
	pack, ans, err := rpc.ReadPack(data, opt)
	if err != nil {
		errorLogger.Println(err.Error())
		return err
	}
	ans, err = service.Pool.Prox(pack)
	if err != nil {
		return err
	}
	rpc.SendAnswer(ans, opt, response)
	return nil
}
