package models

type DelMessage struct {
	Id   string `bson:"_id"`
	Data any    `bson:"data"`
}
