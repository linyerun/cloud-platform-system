package models

// 插入的时候转uint, 目的是保证权限保存类型是Mongo的NumberLong
const (
	LinuxApplicationFormStatusIng    = 0
	LinuxApplicationFormStatusOk     = 1
	LinuxApplicationFormStatusReject = 2
	LinuxApplicationFormDocument     = "linux_application_forms"
)

type LinuxApplicationForm struct { // 这个申请由用户对应的管理员进行审核
	Id          string `bson:"id" json:"id"`
	UserId      string `bson:"user_id" json:"user_id"`
	Explanation string `bson:"explanation" json:"explanation"` // 申请说明

	ImageId       string  `bson:"image_id" json:"image_id"`
	ContainerName string  `bson:"container_name" json:"container_name"`
	ExportPorts   []int64 `bson:"export_ports" json:"export_ports"`

	// 内存相关
	Memory     int64 `bson:"memory" json:"memory"`
	MemorySwap int64 `bson:"memory_swap" json:"memory_swap"` // 允许为-1

	// CPU相关
	CoreCount uint `bson:"core_count" json:"core_count"` // 设置工作线程的数量

	// 磁盘数
	DiskSize uint `bson:"disk_size" json:"disk_size"`

	Status   uint  `bson:"status" json:"status"`
	CreateAt int64 `bson:"create_at" json:"create_at"`
	FinishAt int64 `bson:"finish_at" json:"finish_at"`
}
