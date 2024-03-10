package models

const (
	DbApplicationFormStatusIng    = 0
	DbApplicationFormStatusOk     = 1
	DbApplicationFormStatusReject = 2
	DbApplicationFormDocument     = "db_application_forms"
)

type DbApplicationForm struct {
	Id           string `bson:"_id" json:"id"`
	UserId       string `bson:"user_id" json:"user_id"`
	Explanation  string `bson:"explanation" json:"explanation"`     // 申请说明
	RejectReason string `bson:"reject_reason" json:"reject_reason"` // 拒绝理由

	ImageId string `bson:"image_id" json:"image_id"`
	DbName  string `bson:"db_name" json:"db_name"`

	Status   uint  `bson:"status" json:"status"`
	CreateAt int64 `bson:"create_at" json:"create_at"`
	FinishAt int64 `bson:"finish_at" json:"finish_at"`
}
