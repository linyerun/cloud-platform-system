package models

const LinuxContainerDocument = "linux_containers"

type LinuxPort struct {
	ExportPort uint `bson:"export_port"`
	TargetPort uint `bson:"target_port"`
}

type LinuxContainer struct {
	Id        string `bson:"_id"`
	UserId    string `bson:"user_id"`
	StartTime string `bson:"start_time"`

	Host  string      `bson:"host"`
	Ports []LinuxPort `bson:"ports"`

	Name        string `bson:"name"`
	ContainerId string `bson:"container_id"`
	ImageId     string `bson:"image_id"`

	Status uint `bson:"status"` // 0: 关闭状态, 1: 开启状态

	InitUsername string `bson:"init_username"`
	InitPassword string `bson:"init_password"`
}
