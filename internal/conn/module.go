package conn

import "go.uber.org/fx"

// Module wires the infra connections. Add new connections (RabbitMQ, Temporal,
// Beanstalkd, etc.) by writing a Connect* provider and listing it here.
var Module = fx.Options(
	fx.Provide(
		ConnectPostgres,
		ConnectRedis,
	),
)
