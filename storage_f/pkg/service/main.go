package service

import (
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	"os"
	"storage_f/pkg"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Storage *StorageT

func NewStorage() {
	Storage = &StorageT{
		key: cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE),
	}
	err := Storage.setKey()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	Storage.mail()
}
