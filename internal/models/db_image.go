package models

const (
	DbImageTypeMySql   = "mysql"
	DbImageTypeRedis   = "redis"
	DbImageTypeMongoDb = "mongo"
	DbUsername         = "root"
	DbImageDocument    = "db_images"
)

type DbImage struct {
	Id        string `bson:"_id" json:"id"`
	CreatorId string `bson:"creator_id" json:"creator_id"`

	Type string `bson:"type" json:"type"`

	Name string `bson:"name" json:"name"`
	Tag  string `bson:"tag" json:"tag"`

	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`

	ImageId string `bson:"image_id" json:"image_id"`
	Size    int64  `bson:"size" json:"size"`

	Port uint `bson:"port" json:"port"`

	CreatedAt int64 `bson:"created_at" json:"created_at"`
	IsDeleted bool  `bson:"is_deleted" json:"-"`
}
