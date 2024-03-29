package models

// 插入的时候转uint, 目的是保证权限保存类型是Mongo的NumberLong
const (
	UserDocument   = "users"
	VisitorAuth    = 0
	UserAuth       = 1
	AdminAuth      = 2
	SuperAdminAuth = 3
)

type User struct {
	Id       string `bson:"_id" json:"id,omitempty"`
	Email    string `bson:"email" json:"email,omitempty"`
	Password string `bson:"password" json:"password,omitempty"`
	Name     string `bson:"name" json:"name,omitempty"`
	// Auth不能omitempty, 不然无法录入权限0了
	Auth uint `bson:"auth" json:"auth,omitempty"` // 权限：0: 游客, 1: 用户, 2: 管理员, 3: 超级管理员
}
