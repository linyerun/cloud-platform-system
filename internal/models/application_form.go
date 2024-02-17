package models

// 插入的时候转uint, 目的是保证权限保存类型是Mongo的NumberLong
const (
	ApplicationFormStatusIng    = 0
	ApplicationFormStatusOk     = 1
	ApplicationFormStatusReject = 2
	ApplicationFormTable        = "application_forms"
)

type ApplicationForm struct {
	Id      string `bson:"_id"`
	UserId  string `bson:"user_id"`
	AdminId string `bson:"admin_id"`
	Status  uint   `bson:"status"`
}
