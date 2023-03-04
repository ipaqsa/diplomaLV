package service

import (
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	"os"
	"storage_a/pkg"
	"storage_a/pkg/db"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Storage *StorageT

func NewStorage(path string) {
	Storage = &StorageT{
		key: cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE),
	}
	Storage.db = db.DataBaseInit(path)
	err := Storage.setKey()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	Storage.mail()
}
