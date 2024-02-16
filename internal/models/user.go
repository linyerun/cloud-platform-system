package models

const (
	UserTable      = "users"
	VisitorAuth    = 0
	UserAuth       = 1
	AdminAuth      = 2
	SuperAdminAuth = 3
)

type User struct {
	Id       string `bson:"_id,omitempty"`
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
	Name     string `bson:"name,omitempty"`
	Auth     uint   `bson:"auth,omitempty"` // 权限：0: 游客, 1: 用户, 2: 管理员, 3: 超级管理员
}
