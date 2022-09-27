package server

import (
	"aed-api-server/internal/middleware"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/server/config"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type (
	AdminRouter interface {
		MountAdminRouter(router *route.Router)
	}

	NoAuthRouter interface {
		MountNoAuthRouter(router *route.Router)
	}

	AuthRouter interface {
		MountAuthRouter(router *route.Router)
	}

	OriginEngine interface {
		MountGinEngineRouter(router *route.Router)
	}
)

func initRouter(c *config.AppConfig) {
	eng.Use(middleware.Trace)
	eng.Use(middleware.AccessLog)
	eng.Use(middleware.Recovery)
	eng.Use(middleware.Cors())

	adminGroup := eng.Group("/admin-api")
	noAuthGroup := eng.Group("/api")
	authorizedGroup := eng.Group("/api")
	authorizedGroup.Use(middleware.Authorize)
	adminGroup.Use(middleware.NewAuthorizationAdmin(c.Backend).AuthorizeAdmin)

	route.SetJsonResponseHandlerSucceed(func(ctx *gin.Context, body interface{}) {
		response.ReplyOK(ctx, body)
	})

	route.SetJsonResponseHandlerFailed(func(ctx *gin.Context, body interface{}, err error) {
		response.ReplyErrorWithData(ctx, err, body)
	})

	for _, instance := range component.GetInstances() {
		if adminRouter, ok := instance.Obj.(AdminRouter); ok {
			adminRouter.MountAdminRouter(route.NewRouter(adminGroup))
		}

		if noAuthRouter, ok := instance.Obj.(NoAuthRouter); ok {
			noAuthRouter.MountNoAuthRouter(route.NewRouter(noAuthGroup))
		}

		if authorizedRouter, ok := instance.Obj.(AuthRouter); ok {
			authorizedRouter.MountAuthRouter(route.NewRouter(authorizedGroup))
		}

		if originEng, ok := instance.Obj.(OriginEngine); ok {
			originEng.MountGinEngineRouter(route.NewRouter(eng))
		}
	}
}
