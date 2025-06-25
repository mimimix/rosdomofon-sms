package app

import (
	webServer "domofon-api/internal/transport/http"
	httpHandlers "domofon-api/internal/transport/http/handler"
	"domofon-api/pkg/rosdomofon"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		webServer.New,
		rosdomofon.NewDomofon,
	),
	httpHandlers.HttpHandlers,
)
