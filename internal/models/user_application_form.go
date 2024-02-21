package models

// 插入的时候转uint, 目的是保证权限保存类型是Mongo的NumberLong
const (
	UserApplicationFormStatusIng    = 0
	UserApplicationFormStatusOk     = 1
	UserApplicationFormStatusReject = 2
	UserApplicationFormDocument     = "user_application_forms"
)

type UserApplicationForm struct {
	Id          string `bson:"_id"`
	UserId      string `bson:"user_id"`
	AdminId     string `bson:"admin_id"`
	Explanation string `bson:"explanation"`
	Status      uint   `bson:"status"`
}
