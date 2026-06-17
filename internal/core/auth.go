package core

import (
	"time"

	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/cache"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/internal/stores"
)

// Auth is the domain use-case object for the authentication flows (signup
// initiate + callback). It owns the shared dependencies and builds the base
// flow Ctx, mirroring Widget.
type Auth struct {
	cfg   *config.Config
	srv   *services.SrvHolder
	store *stores.StoHolder
	cache cache.Cache
	log   *zap.Logger
}

func NewAuth(
	cfg *config.Config,
	srv *services.SrvHolder,
	store *stores.StoHolder,
	cacheStore cache.Cache,
	log *zap.Logger,
) *Auth {
	return &Auth{
		cfg:   cfg,
		srv:   srv,
		store: store,
		cache: cacheStore,
		log:   log,
	}
}

// BaseCtx builds a flow context pre-populated with shared dependencies.
func (a *Auth) BaseCtx(lang string) Ctx {
	return Ctx{
		Srv:    a.srv,
		Store:  a.store,
		Cache:  a.cache,
		Config: a.cfg,
		Log:    a.log,
		Now:    time.Now(),
		Lang:   lang,
	}
}

// TransformUser builds the API response from a stored user.
func (a *Auth) TransformUser(u *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            u.ID,
		WorkOSUserID:  u.WorkOSUserID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		EmailVerified: u.EmailVerified,
		Status:        u.Status,
		CreatedAt:     u.CreatedAt.Unix(),
	}
}
