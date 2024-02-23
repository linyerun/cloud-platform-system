package models

type DbContainer struct {
	Id       string `bson:"_id"`
	Type     string `bson:"type"`
	CreateAt int64  `bson:"create_at"`
}
