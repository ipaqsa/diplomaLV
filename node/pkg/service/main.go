package service

import (
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	"node/pkg"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

var Node *NodeT

func NewNode() {
	Node = &NodeT{
		Key: cryptoUtils.GeneratePrivate(pkg.Config.AKEY_SIZE),
	}
	//err := Node.setKey()
	//if err != nil {
	//	println(err.Error())
	//	os.Exit(-1)
	//}
}
