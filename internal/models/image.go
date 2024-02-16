package models

type Image struct {
	Id        string `bson:"_id,omitempty"`
	CreatorId string `bson:"creator_id,omitempty"`
}
