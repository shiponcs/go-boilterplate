package stores

import (
	"gorm.io/gorm"

	"github.com/example/go-svc-boilerplate/internal/models/entity"
)

// WidgetStore is a GORM-backed repository for widgets. Stores hold no business
// logic — just data access. Mirror this shape for each new entity.
type WidgetStore struct {
	db *gorm.DB
}

func NewWidgetStore(db *gorm.DB) *WidgetStore {
	return &WidgetStore{db: db}
}

func (s *WidgetStore) GetByID(id uint) (*entity.Widget, error) {
	var w entity.Widget
	if err := s.db.First(&w, id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *WidgetStore) Create(w *entity.Widget) error {
	return s.db.Create(w).Error
}

func (s *WidgetStore) List(limit int) ([]entity.Widget, error) {
	var widgets []entity.Widget
	if err := s.db.Limit(limit).Order("id desc").Find(&widgets).Error; err != nil {
		return nil, err
	}
	return widgets, nil
}
