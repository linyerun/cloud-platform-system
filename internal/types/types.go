// Code generated by goctl. DO NOT EDIT.
package types

type ApplicationFormPostRequest struct {
	AdminId    string `json:"admin_id"`
	AdminEmail string `json:"admin_email"`
}

type CaptchaEmailRequest struct {
	Email string `form:"email"`
}

type CaptchaPictureResponse struct {
	PicData []byte `json:"pic_data"`
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

type DeleteContainerRequest struct {
	ContainerId string `path:"id"`
}

type DeleteUserRequest struct {
	UserId string `path:"id"`
}

type GetApplicationFormByStatusRequest struct {
	Status uint `path:"status"`
}

type GetFormByStatusRequest struct {
	Status uint `path:"status"`
}

type ImageDelRequest struct {
	ImageName string `json:"image_name"`
}

type ImagePullRequest struct {
	ImageName string `json:"image_name"`
	ImageTag  string `json:"image_tag"`
}

type PutVisitorToUserRequest struct {
	VisitorId    string `path:"visitor_id"`
	Status       uint   `path:"status"`
	VisitorEmail string `json:"visitor_email"`
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
