package models

const DelMessageDocument = "del_message"

type DelMessage struct {
	Id   string `bson:"_id"`
	Data any    `bson:"data"`
}

func NewDelMessage(id string, data any) *DelMessage {
	return &DelMessage{Id: id, Data: data}
}
