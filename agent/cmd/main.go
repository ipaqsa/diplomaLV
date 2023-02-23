package main

import (
	"agent/pkg"
	"agent/pkg/server"
	"github.com/ipaqsa/netcom/configurator"
	"os"
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
	err := server.Run()
	if err != nil {
		println(err.Error())
		return
	}
}
