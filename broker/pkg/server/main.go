package server

import (
	"broker/pkg"
	"broker/pkg/service"
	"github.com/ipaqsa/netcom/logger"
	netcom "github.com/ipaqsa/netcom/rpc"
	rpc "net/rpc"
)

type ServerBroker int

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var server ServerBroker

func Run() error {
	service.NewBroker()

	server = 0
	err := rpc.Register(&server)

	if err != nil {
		return err
	}
	infoLogger.Printf("broker starting on %s\n", pkg.Config.Port)
	err = netcom.ListenRPC(pkg.Config.Port)
	if err != nil {
		return err
	}
	return nil
}
