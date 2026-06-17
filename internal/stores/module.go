package stores

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewWidgetStore,
		NewUserStore,
		NewStoHolder,
	),
)
