package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewHttpClient,
		fx.Annotate(
			NewPricingService,
			fx.As(new(PricingService)),
		),
		fx.Annotate(
			NewWorkOSAuthService,
			fx.As(new(AuthService)),
		),
		NewSrvHolder,
	),
)
