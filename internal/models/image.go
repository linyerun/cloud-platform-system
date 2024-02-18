package models

const ImageDocument = "images"

type Image struct {
	Id        string `bson:"_id"`
	CreatorId string `bson:"creator_id"`
	Name      string `bson:"name"`
	Tag       string `bson:"tag"`
	ImageId   string `bson:"image_id"`
	Size      int64  `bson:"size"`
}
