package cache

import "go.uber.org/fx"

// Module provides the concrete *store as both the Store and Cache interfaces,
// so consumers can depend on the narrowest surface they need.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewStore,
			fx.As(new(Store)),
			fx.As(new(Cache)),
		),
	),
)
