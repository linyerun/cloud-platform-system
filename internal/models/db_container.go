package models

const (
	DbContainerDocument       = "db_containers"
	DbContainerStatusRunning  = 0
	DbContainerStatusSleeping = 1
	DbContainerDel            = 3
)

type DbContainer struct {
	Id              string `bson:"_id" json:"id"`                              // ID
	UserId          string `bson:"user_id" json:"user_id"`                     // 容器的所有者
	Name            string `bson:"name" json:"name"`                           // 容器名称, 打算让它等于雪花ID使其唯一, 用它就行了
	DbContainerName string `bson:"db_container_name" json:"db_container_name"` // 用户为它取的名字
	ImageId         string `bson:"image_id" json:"image_id"`                   // 指: linux_images文档的_id属性

	CreateAt  int64 `bson:"create_at" json:"create_at"`   // 创建时间
	StartTime int64 `bson:"start_time" json:"start_time"` // 启动时间
	StopTime  int64 `bson:"stop_time" json:"stop_time"`   // 关闭时间
	Status    uint  `bson:"status" json:"status"`         // 0: 关闭状态, 1: 开启状态

	Host string `bson:"host" json:"host"` // 主机地址
	Port uint   `bson:"port"`             // 实际端口

	Type     string `bson:"type" json:"type"`         // 数据库类型
	Username string `bson:"username" json:"username"` // 登录账号
	Password string `bson:"password" json:"password"` // 登录密码
}
