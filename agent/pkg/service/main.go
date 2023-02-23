package service

import (
	"agent/pkg"
	"agent/pkg/db"
	"crypto/rsa"
	"github.com/ipaqsa/netcom/logger"
	"os"
)

var tasks []Task

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Service *Agent

func NewAgent() {
	Service = &Agent{
		storage: Storage{
			data: make(map[string]*rsa.PrivateKey),
		},
	}

	infoLogger.Printf("create agent")
	runService()
}

func runService() {
	err := db.InitDB(pkg.Config.DBPath)
	if err != nil && err.Error() != "create bucket: bucket already exists" {
		errorLogger.Printf("%s", err.Error())
		os.Exit(-1)
	}
	err = Service.load()
	if err != nil {
		errorLogger.Printf("%s", err.Error())
		os.Exit(-1)
	}
	for _, service := range pkg.Config.Services {
		err = Service.save(service)
		if err != nil && err.Error() != "key exist" {
			errorLogger.Printf("%s", err.Error())
			os.Exit(-1)
		}
	}
	Service.Key = Service.Get("agent")
	if Service.Key == nil {
		errorLogger.Printf("selfkey is nil")
		os.Exit(-1)
	}
	go Service.mail()
	go Service.fetch()
}
