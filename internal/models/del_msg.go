package models

type DelMessage struct {
	Id   string `bson:"_id,omitempty"`
	Data any    `bson:"data,omitempty"`
}
