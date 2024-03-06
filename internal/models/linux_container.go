package models

const (
	LinuxSleep             = 0
	LinuxRunning           = 1
	LinuxContainerDocument = "linux_containers"
)

type LinuxContainer struct {
	Id                string `bson:"_id" json:"id"`
	UserId            string `bson:"user_id" json:"user_id"`                         // 容器的所有者
	Name              string `bson:"name" json:"name"`                               // 容器名称, 打算让它等于雪花ID使其唯一, 用它就行了
	UserContainerName string `bson:"user_container_name" json:"user_container_name"` // 用户为它取的名字
	ImageId           string `bson:"image_id"`                                       // 指: linux_images文档的_id属性

	CreateAt  int64 `bson:"create_at" json:"create_at"`   // 容器创建时间
	StartTime int64 `bson:"start_time" json:"start_time"` // 容器启动时间
	StopTime  int64 `bson:"stop_time" json:"stop_time"`   // 容器关闭时间
	Status    uint  `bson:"status" json:"status"`         // 0: 关闭状态, 1: 开启状态

	Host         string          `bson:"host" json:"host"`                   // 主机地址
	PortsMapping map[int64]int64 `bson:"ports_mapping" json:"ports_mapping"` // 端口映射

	// 初始化容器的账户与密码
	InitUsername string `bson:"init_username" json:"init_username"`
	InitPassword string `bson:"init_password" json:"init_password"`

	// 内存相关(单位: 字节)
	Memory     int64 `bson:"memory" json:"memory"`
	MemorySwap int64 `bson:"memory_swap" json:"memory_swap"`

	// CPU相关
	CoreCount uint `bson:"core_count" json:"core_count"` // 设置工作线程的数量

	// 磁盘数(单位: 字节)
	DiskSize uint `bson:"disk_size" json:"disk_size"`
}
