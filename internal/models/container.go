package models

type Container struct {
	Id     string `bson:"_id,omitempty"`
	UserId string `bson:"user_id,omitempty"`
}
