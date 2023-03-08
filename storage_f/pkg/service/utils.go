package service

import (
	"encoding/json"
	"errors"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"os"
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

func isExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func UnmarshalMessage(data string) (*FileMessage, error) {
	var file FileMessage
	err := json.Unmarshal([]byte(data), &file)
	if err != nil {
		return nil, err
	}
	return &file, err
}

func MarshalFile(file []byte, filename, sender string) (string, error) {
	data := cryptoUtils.Base64Encode(file)
	fm := FileMessage{Title: filename, Data: data, Meta: sender}
	jfm, err := json.Marshal(fm)
	return string(jfm), err
}
