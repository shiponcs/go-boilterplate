package cache

import (
	"time"

	"github.com/example/go-svc-boilerplate/internal/models/entity"
)

// Store is the read/write cache surface used by core use-cases. Split larger
// surfaces into multiple interfaces (and compose them into Cache) as the
// service grows — this keeps consumers depending only on what they use.
type Store interface {
	GetWidget(id uint) (*entity.Widget, error)
	SetWidget(id uint, widget *entity.Widget, ttl time.Duration) error
	ForgetWidget(id uint) error
}

// Cache is the full cache surface. NewStore satisfies it; fx provides the
// concrete type as both Store and Cache (see module.go).
type Cache interface {
	Store
}
