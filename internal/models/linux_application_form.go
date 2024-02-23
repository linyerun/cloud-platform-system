package models

// 插入的时候转uint, 目的是保证权限保存类型是Mongo的NumberLong
const (
	LinuxApplicationFormStatusIng    = 0
	LinuxApplicationFormStatusOk     = 1
	LinuxApplicationFormStatusReject = 2
	LinuxApplicationFormDocument     = "linux_application_forms"
)

type LinuxApplicationForm struct { // 这个申请由用户对应的管理员进行审核
	UserId      string  `bson:"user_id"`
	ImageId     string  `bson:"image_id"`
	ExportPorts []int64 `bson:"export_ports"`
	Explanation string  `bson:"explanation"` // 申请说明
	Status      uint    `bson:"status"`
	CreateAt    int64   `bson:"create_at"`
	FinishAt    int64   `bson:"finish_at"`
}
