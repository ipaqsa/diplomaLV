package db

import "encoding/json"

func MarshalMessages(msg *Messages) string {
	jmsg, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(jmsg)
}

func UnmarshalMessage(data string) *Message {
	var msg Message
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		return nil
	}
	return &msg
}
