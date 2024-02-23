package models

const (
	AsyncTaskIng      = 0
	AsyncTaskOk       = 1
	AsyncTaskFail     = 2
	AsyncTaskDocument = "async_tasks"
)

type AsyncTask struct {
	Id       string `bson:"_id" json:"id"`
	UserId   string `bson:"user_id" json:"user_id"`
	Type     string `bson:"type" json:"type"`
	Args     any    `bson:"args" json:"args"`
	RespData any    `bson:"resp_data" json:"resp_data,omitempty"`
	Priority uint   `bson:"priority" json:"priority"` // 优先级, 数字越搞优先级越高
	Status   uint   `bson:"status" json:"status"`
	CreateAt int64  `bson:"create_at" json:"create_at"`
	FinishAt int64  `bson:"finish_at" json:"finish_at,omitempty"`
}
