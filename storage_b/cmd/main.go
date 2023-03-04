package main

import (
	"github.com/ipaqsa/netcom/configurator"
	"os"
	"storage_b/pkg"
	"storage_b/pkg/service"
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
	//service.NewStorage(pkg.Config.PathToDB)
	service.NewStorage("data/msg.db")
}
