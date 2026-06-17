package conn

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/example/go-svc-boilerplate/internal/config"
)

func ConnectPostgres(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	pg := cfg.Postgres
	db, err := gorm.Open(postgres.Open(pg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(pg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(pg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(pg.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("successfully connected to database")
	return db, nil
}
