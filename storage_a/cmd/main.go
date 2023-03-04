package main

import (
	"github.com/ipaqsa/netcom/configurator"
	"os"
	"storage_a/pkg"
	"storage_a/pkg/service"
)

func init() {
	err := configurator.InitConfiguration(&pkg.Config, "0.0.1")
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	configurator.InitInfo(pkg.Config.Port)
}

func main() {
	service.NewStorage("data/storagea.db")
}
