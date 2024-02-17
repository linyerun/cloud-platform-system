package models

type Image struct {
	Id        string `bson:"_id"`
	CreatorId string `bson:"creator_id"`
}
