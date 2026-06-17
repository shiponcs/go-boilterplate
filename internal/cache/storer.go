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

	// SetSignupState records a one-time AuthKit CSRF state value with a TTL.
	SetSignupState(state string, ttl time.Duration) error
	// ConsumeSignupState atomically checks-and-deletes a state value, returning
	// true if it was present (i.e. issued by us and not yet used).
	ConsumeSignupState(state string) (bool, error)
}

// Cache is the full cache surface. NewStore satisfies it; fx provides the
// concrete type as both Store and Cache (see module.go).
type Cache interface {
	Store
}
