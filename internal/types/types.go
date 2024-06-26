// Code generated by goctl. DO NOT EDIT.
package types

type ApplicationFormPostRequest struct {
	AdminId     string `json:"admin_id"`
	AdminEmail  string `json:"admin_email"`
	Explanation string `json:"explanation"`
}

type CaptchaEmailRequest struct {
	Email string `form:"email"`
}

type CaptchaPictureResponse struct {
	PicData []byte `json:"pic_data"`
}

type ChangeDbApplicationStatusReq struct {
	Id           string `path:"id"`
	Status       uint   `path:"status"`
	RejectReason string `json:"reject_reason"`
}

type ChangeForgetPasswordReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type ChangeUserMsgReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type CreateAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type DbStartApplyReq struct {
	DbName      string `json:"db_name"`
	ImageId     string `json:"image_id"`
	Explanation string `json:"explanation"`
}

type DelDbImageByIdReq struct {
	Id string `path:"id"`
}

type DelExceptionByIdxReq struct {
	Idx uint `path:"idx"`
}

type DelLinuxStopContainerReq struct {
	ContainerId string `path:"container_id"`
}

type DeleteContainerRequest struct {
	ContainerId string `path:"id"`
}

type DeleteUserRequest struct {
	UserId string `path:"id"`
}

type GetAdminMsgByIdReq struct {
	Id string `path:"id"`
}

type GetApplicationFormByStatusRequest struct {
	Status uint `path:"status"`
}

type GetDbImageByIdReq struct {
	Id string `path:"id"`
}

type GetFormByStatusRequest struct {
}

type GetImageMsgByIdReq struct {
	Id string `path:"id"`
}

type GetImageMsgByIdResp struct {
	CreatorName     string   `json:"creator_name"`
	CreatorEmail    string   `json:"creator_email"`
	ImageName       string   `json:"image_name"`
	ImageTag        string   `json:"image_tag"`
	ImageSize       int64    `json:"image_size"`
	EnableCommands  []string `json:"enable_commands"`
	MustExportPorts []int64  `json:"must_export_ports"`
}

type GetUserMsgByIdReq struct {
	Id string `path:"id"`
}

type GetUserMsgByIdResp struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GetUserMsgResp struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type HandleUserLinuxApplicationReq struct {
	FormId string `path:"form_id"`
	Status uint   `path:"status"`
}

type ImageDelRequest struct {
	Id string `path:"id"`
}

type ImagePullRequest struct {
	ImageName            string   `json:"image_name"`
	ImageTag             string   `json:"image_tag"`
	ImageEnabledCommands []string `json:"image_enabled_commands"`
	ImageMustExportPorts []int64  `json:"image_must_export_ports"`
}

type LinuxStartApplyRequest struct {
	ContainerName string  `json:"container_name"`
	ImageId       string  `json:"image_id"`
	ExportPorts   []int64 `json:"export_ports"`
	Explanation   string  `json:"explanation"`
	Memory        int64   `json:"memory"`
	MemorySwap    int64   `json:"memory_swap"`
	CoreCount     uint    `json:"core_count"`
	DiskSize      uint    `json:"disk_size"`
}

type PullDbImageReq struct {
	ImageName string `json:"image_name"`
	ImageTag  string `json:"image_tag"`
	Type      string `json:"type"`
	Port      uint   `json:"port"`
}

type PutVisitorToUserRequest struct {
	VisitorId    string `path:"visitor_id"`
	Status       uint   `path:"status"`
	VisitorEmail string `json:"visitor_email"`
}

type UpdateDbStatusReq struct {
	DbId   string `path:"db_id"`
	Status uint   `path:"status"`
}

type UpdateLinuxStatusReq struct {
	ContainerId string `path:"container_id"`
	Status      uint   `path:"status"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Captcha  string `json:"captcha"`
}
