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
					Path:    "/user/image/list",
					Handler: v2.UserGetImagesHandler(serverCtx),
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
					Path:    "/admin/image/list",
					Handler: v3.AdminGetImagesHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/application_forms/:status",
					Handler: v3.GetFormByStatusHandler(serverCtx),
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
					Method:  http.MethodDelete,
					Path:    "/image/del/:id",
					Handler: v3.DeleteImageHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/image/pull",
					Handler: v3.PullImageHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/user/del/:id",
					Handler: v3.DeleteUserByIdHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users",
					Handler: v3.GetUsersHandler(serverCtx),
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
					Path:    "/token/refresh",
					Handler: v5.RefreshTokenHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/v5"),
	)
}
