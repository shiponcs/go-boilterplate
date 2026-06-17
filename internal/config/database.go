package config

import "fmt"

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DB              string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime int // seconds
}

// GetDSN returns the Postgres Data Source Name for the connection.
func (dc *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dc.Host,
		dc.User,
		dc.Password,
		dc.DB,
		dc.Port,
	)
}
