package models

const (
	ApplicationFormStatusIng    = 0
	ApplicationFormStatusOk     = 1
	ApplicationFormStatusReject = 2
	ApplicationFormTable        = "application_forms"
)

type ApplicationForm struct {
	Id      string `bson:"_id,omitempty"`
	UserId  string `bson:"user_id"`
	AdminId string `bson:"admin_id"`
	Status  uint   `bson:"status"`
}
