package app

import (
	"domofon-api/connections/modem"
	checker "domofon-api/internal"
	"domofon-api/pkg/smsPoller"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		modem.New,
		smsPoller.New,
	),
	fx.Invoke(
		checker.Start,
	),
)
