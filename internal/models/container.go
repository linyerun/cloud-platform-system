package models

type Container struct {
	Id     string `bson:"_id"`
	UserId string `bson:"user_id"`
}
