// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	v0 "cloud-platform-system/internal/handler/v0"
	v1 "cloud-platform-system/internal/handler/v1"
	v2 "cloud-platform-system/internal/handler/v2"
	v3 "cloud-platform-system/internal/handler/v3"
	v4 "cloud-platform-system/internal/handler/v4"
	v5 "cloud-platform-system/internal/handler/v5"
	"cloud-platform-system/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/captcha",
				Handler: v0.CaptchaPictureHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/captcha/email",
				Handler: v0.CaptchaEmailHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/user/change/pwd",
				Handler: v0.ChangeForgetPasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: v0.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: v0.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v0"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth, serverCtx.Visitor},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/admin/all",
					Handler: v1.GetAdminsHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/to/user",
					Handler: v1.ToUserHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth, serverCtx.User},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/admin/:id",
					Handler: v2.GetAdminMsgByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/db-container/change/:db_id/:status",
					Handler: v2.UpdateDbStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/db-container/list",
					Handler: v2.GetDbContainerHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/linux_application_form/list",
					Handler: v2.GetLinuxApplicationFormListHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/linux_container/:container_id",
					Handler: v2.DelLinuxStopContainerHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/linux_container/change/:container_id/:status",
					Handler: v2.UpdateLinuxStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/linux_container/list",
					Handler: v2.GetLinuxContainerByUserIdHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/myself/db-application/list",
					Handler: v2.GetDbApplicationByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/apply/db",
					Handler: v2.DbStartApplyHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/apply/linux",
					Handler: v2.LinuxStartApplyHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/db-image/list",
					Handler: v2.UserGetDbImagesHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/linux/list",
					Handler: v2.UserGetLinuxImagesHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v2"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth, serverCtx.Admin},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/admin/linux/list",
					Handler: v3.AdminGetLinuxImagesHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/application_forms/list",
					Handler: v3.GetFormByStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/asynctask/list",
					Handler: v3.GetAsyncTaskListHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/container/:id",
					Handler: v3.DeleteContainerHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/container/list",
					Handler: v3.GetContainerListHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/db-application/:id/:status",
					Handler: v3.ChangeDbApplicationStatusHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/db-application/list",
					Handler: v3.GetDbApplicationListHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/db-image/:id",
					Handler: v3.DelDbImageByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/db-image/:id",
					Handler: v3.GetDbImageByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/db-image/list",
					Handler: v3.GetDbImageListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/db-image/pull",
					Handler: v3.PullDbImageHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/image/:id",
					Handler: v3.GetImageMsgByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/linux/del/:id",
					Handler: v3.DeleteLinuxImageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/linux/pull",
					Handler: v3.PullLinuxImageHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/:id",
					Handler: v3.GetUserMsgByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/user/del/:id",
					Handler: v3.DeleteUserByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/user_linux_apply/handle/:form_id/:status",
					Handler: v3.HandleUserLinuxApplicationHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users",
					Handler: v3.GetUsersHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users_linux_application_form/list",
					Handler: v3.GetUsersLinuxApplicationFormListHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/visitor_to_user/:visitor_id/:status",
					Handler: v3.VisitorToUserHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v3"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth, serverCtx.Super},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/create_admin",
					Handler: v4.CreateAdminHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/exception/:idx",
					Handler: v4.DelExceptionByIdxHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/exception/list",
					Handler: v4.GetExceptionListHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v4"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/change/user-msg",
					Handler: v5.ChangeUserMsgHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/token/refresh",
					Handler: v5.RefreshTokenHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user-msg",
					Handler: v5.GetUserMsgHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v5"),
	)
}
