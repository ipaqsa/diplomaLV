package service

import (
	"admin/pkg"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	"os"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Admin *AdminT

func NewAdmin() {
	Admin = &AdminT{
		key: cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE),
	}
	err := Admin.setKey()
	infoLogger.Println("key`s set")
	if err != nil {
		errorLogger.Println(err.Error())
		os.Exit(-1)
	}
	infoLogger.Println("admin`s created")
	Admin.mail()
}
