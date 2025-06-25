package httpHandlers

import (
	"domofon-api/internal/transport/http/handler/ApiRouters"
	"domofon-api/internal/transport/http/handler/api"

	"go.uber.org/fx"
)

var HttpHandlers = fx.Module("httpHandlers",
	fx.Provide(
		ApiRouters.CreateApiRoutes,
		fx.Private,
	),
	fx.Invoke(
		apiRoute.ApiRoute,
	),
)
