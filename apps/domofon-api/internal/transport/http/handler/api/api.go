package apiRoute

import (
	"domofon-api/internal/transport/http/handler/ApiRouters"
	"domofon-api/pkg/rosdomofon"

	"domofon-api.gg/config"

	"go.uber.org/fx"
)

type Route struct {
	routers    *ApiRouters.ApiRouters
	config     *config.Config
	rosdomofon *rosdomofon.Domofon
}

type fxOpts struct {
	fx.In
	ApiRouter  *ApiRouters.ApiRouters
	Config     *config.Config
	Rosdomofon *rosdomofon.Domofon
}

func ApiRoute(opts fxOpts) *Route {
	router := &Route{
		routers:    opts.ApiRouter,
		config:     opts.Config,
		rosdomofon: opts.Rosdomofon,
	}

	opts.ApiRouter.Public.GET("/open", router.open)

	return router
}
