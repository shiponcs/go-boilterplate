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
	WorkOS   WorkOSConfig
	Session  SessionConfig
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
		WorkOS: WorkOSConfig{
			ApiKey:      v.GetString("workos.api_key"),
			ClientID:    v.GetString("workos.client_id"),
			RedirectURI: v.GetString("workos.redirect_uri"),
		},
		Session: SessionConfig{
			CookieName:     v.GetString("session.cookie_name"),
			CookiePassword: v.GetString("session.cookie_password"),
			CookieDomain:   v.GetString("session.cookie_domain"),
			Secure:         v.GetBool("session.secure"),
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

	// WorkOS credentials have no safe defaults; set them via env.
	v.SetDefault("workos.api_key", "")
	v.SetDefault("workos.client_id", "")
	v.SetDefault("workos.redirect_uri", "http://localhost:8080/api/v1/auth/callback")

	// Session cookie. cookie_password seals the cookie and must be a 32-byte
	// key set via env; cookie_domain is optional.
	v.SetDefault("session.cookie_name", "wos_session")
	v.SetDefault("session.cookie_password", "")
	v.SetDefault("session.cookie_domain", "")
	v.SetDefault("session.secure", true)
}
