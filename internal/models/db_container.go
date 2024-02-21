package models

type DbContainer struct {
	Id   string `bson:"_id"`
	Type string `bson:"type"`
}
