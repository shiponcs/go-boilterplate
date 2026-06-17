package stores

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/example/go-svc-boilerplate/internal/models/entity"
)

// UserStore is a GORM-backed repository for the local user mirror. Stores hold
// no business logic — just data access.
type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetByWorkOSID(workosUserID string) (*entity.User, error) {
	var u entity.User
	if err := s.db.Where("workos_user_id = ?", workosUserID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Upsert inserts the user or, if one already exists for the same
// workos_user_id, updates the mutable profile fields. The passed entity is
// back-filled with the stored row (including ID/CreatedAt).
func (s *UserStore) Upsert(u *entity.User) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "workos_user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"email", "first_name", "last_name", "email_verified", "status", "updated_at"}),
	}).Create(u).Error
}
