package server

import (
	"github.com/ipaqsa/netcom/logger"
	netcom "github.com/ipaqsa/netcom/rpc"
	"ingress/pkg"
	"ingress/pkg/service"
	"net/rpc"
)

type IngressT int

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var ingress IngressT

func Run() error {
	service.CreatePool()
	ingress = 0
	err := rpc.Register(&ingress)

	if err != nil {
		return err
	}
	infoLogger.Printf("ingress starting on %s\n", pkg.Config.Port)
	err = netcom.ListenRPC(pkg.Config.Port)
	if err != nil {
		return err
	}
	return nil
}
