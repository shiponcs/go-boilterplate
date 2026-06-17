package config

import (
	"strings"

	"github.com/spf13/viper"
	// Uncomment to enable the Consul remote provider (see LoadConfig note).
	// _ "github.com/spf13/viper/remote"
)

// Config is the fully-resolved configuration for the service. Sub-structs live
// in their own files (app.go, database.go, redis.go, ...) so each concern is
// easy to find and extend.
type Config struct {
	App      AppConfig
	Postgres DatabaseConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Pricing  PricingConfig
}

// LoadConfig reads configuration from environment variables with sane defaults.
//
// Env keys mirror the struct paths with "." replaced by "_": app.port -> APP_PORT,
// postgres.host -> POSTGRES_HOST, etc.
//
// To back configuration with Consul instead (as is common in production), bind a
// consul URL/path and use viper's remote provider before reading values:
//
//	v.BindEnv("consul_url")
//	v.BindEnv("consul_path")
//	_ = v.AddRemoteProvider("consul", v.GetString("consul_url"), v.GetString("consul_path"))
//	v.SetConfigType("yaml")
//	_ = v.ReadRemoteConfig()
//
// (requires the `_ "github.com/spf13/viper/remote"` import above).
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaults(v)

	cfg := &Config{
		App: AppConfig{
			Port: v.GetString("app.port"),
			Env:  v.GetString("app.env"),
		},
		Postgres: DatabaseConfig{
			Host:            v.GetString("postgres.host"),
			Port:            v.GetInt("postgres.port"),
			User:            v.GetString("postgres.user"),
			Password:        v.GetString("postgres.password"),
			DB:              v.GetString("postgres.db"),
			MaxIdleConn:     v.GetInt("postgres.max_idle_conn"),
			MaxOpenConn:     v.GetInt("postgres.max_open_conn"),
			ConnMaxLifetime: v.GetInt("postgres.conn_max_lifetime"),
		},
		Redis: RedisConfig{
			Host: v.GetString("redis.host"),
			Port: v.GetInt("redis.port"),
			DB:   v.GetInt("redis.db"),
		},
		Logger: LoggerConfig{
			Level: v.GetString("logger.level"),
		},
		Pricing: PricingConfig{
			BaseFare: v.GetFloat64("pricing.base_fare"),
			PerUnit:  v.GetFloat64("pricing.per_unit"),
		},
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.env", "local")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.db", "boilerplate")
	v.SetDefault("postgres.max_idle_conn", 10)
	v.SetDefault("postgres.max_open_conn", 100)
	v.SetDefault("postgres.conn_max_lifetime", 300)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	v.SetDefault("logger.level", "info")

	v.SetDefault("pricing.base_fare", 50.0)
	v.SetDefault("pricing.per_unit", 10.0)
}
