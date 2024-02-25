package models

import "time"

const DelMessageDocument = "del_message"

type DelMessage struct {
	Id        string `bson:"_id"`
	Data      any    `bson:"data"`
	DeletedAt int64  `bson:"deleted_at"`
	Document  string `bson:"document"`
}

func NewDelMessage(id string, data any, document string) *DelMessage {
	return &DelMessage{Id: id, Data: data, DeletedAt: time.Now().UnixMilli(), Document: document}
}
