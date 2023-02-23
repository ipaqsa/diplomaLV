package service

import (
	"broker/pkg"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	"github.com/ipaqsa/netcom/packUtils"
	"os"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Broker *BrokerT

func NewBroker() {
	broker := BrokerT{
		cache: Cache{
			data: make(map[string]Task),
		},
		queues: make(map[string]*Queue),
		reports: Reports{
			data: make(map[string]*packUtils.Package),
		},
		key: cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE),
	}
	err := broker.setKey()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	broker.initServices()
	go broker.cacheMonitor()
	Broker = &broker
}
