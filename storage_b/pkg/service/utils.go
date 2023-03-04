package service

import (
	"github.com/ipaqsa/netcom/cryptoUtils"
	"strings"
)

func hashPrepare(sender, receiver string) (string, string) {
	sender = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(sender))), "/", "S")
	receiver = strings.ReplaceAll(cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(receiver))), "/", "S")

	sender = strings.ReplaceAll(sender, "+", "S")
	receiver = strings.ReplaceAll(receiver, "+", "S")

	sender = strings.ReplaceAll(sender, "=", "S")
	receiver = strings.ReplaceAll(receiver, "=", "S")

	sender = strings.ReplaceAll(sender, "?", "S")
	receiver = strings.ReplaceAll(receiver, "?", "S")

	return sender, receiver
}
