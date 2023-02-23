package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

func SHA1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func taskMarshal(task *Task) string {
	marshal, _ := json.Marshal(task)
	return string(marshal)
}
