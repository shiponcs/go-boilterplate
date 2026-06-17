package conn

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/example/go-svc-boilerplate/internal/models/entity"
)

// RunMigrations creates/updates schema for entities owned by this service.
// Wired as an fx.Invoke so it runs once at startup. The widgets table is
// assumed to be provisioned externally; only add entities here that this
// service is the source of truth for.
func RunMigrations(db *gorm.DB, log *zap.Logger) error {
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}
	log.Info("database migrations applied")
	return nil
}
