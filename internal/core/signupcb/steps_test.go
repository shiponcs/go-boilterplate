package signupcb

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// fakeCache implements cache.Cache; only the signup-state methods are
// exercised. stateOK controls whether ConsumeSignupState reports a hit.
type fakeCache struct {
	stateOK bool
}

func (f *fakeCache) GetWidget(uint) (*entity.Widget, error)              { return nil, nil }
func (f *fakeCache) SetWidget(uint, *entity.Widget, time.Duration) error { return nil }
func (f *fakeCache) ForgetWidget(uint) error                             { return nil }
func (f *fakeCache) SetSignupState(string, time.Duration) error          { return nil }
func (f *fakeCache) ConsumeSignupState(string) (bool, error)             { return f.stateOK, nil }

func newCtx(cache *fakeCache, code, state string) *SignupCbCtx {
	ctx := &SignupCbCtx{Code: code, State: state}
	ctx.Cache = cache
	return ctx
}

func assertBadReq(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var sc errs.StatusCoder
	if !errors.As(err, &sc) {
		t.Fatalf("error %v does not carry a status code", err)
	}
	if sc.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", sc.StatusCode())
	}
}

func TestValidate_MissingCode(t *testing.T) {
	v := &validate{ctx: newCtx(&fakeCache{stateOK: true}, "  ", "state")}
	assertBadReq(t, v.Do(&core.DoCtx{}))
}

func TestValidate_StateMismatch(t *testing.T) {
	v := &validate{ctx: newCtx(&fakeCache{stateOK: false}, "code", "bad-state")}
	assertBadReq(t, v.Do(&core.DoCtx{}))
}

func TestValidate_OK(t *testing.T) {
	v := &validate{ctx: newCtx(&fakeCache{stateOK: true}, "code", "good-state")}
	if err := v.Do(&core.DoCtx{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
