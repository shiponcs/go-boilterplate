package authorize

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/value"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// fakeAuth implements services.AuthService. ValidateAccessToken reports the
// "expired" token as expired and treats "fresh" as valid; RefreshSession hands
// back a "fresh" access token.
type fakeAuth struct {
	refreshErr error
}

func (fakeAuth) SignupURL(string) (string, error)         { return "", nil }
func (fakeAuth) LogoutURL(string, string) (string, error) { return "", nil }
func (fakeAuth) AuthenticateWithCode(context.Context, string) (*value.AuthnResult, error) {
	return nil, nil
}
func (fakeAuth) SessionIDFromToken(string) (string, error) { return "sess_1", nil }

func (fakeAuth) ValidateAccessToken(token string) (*value.SessionClaims, error) {
	switch token {
	case "valid", "fresh":
		return &value.SessionClaims{WorkOSUserID: "user_1", WorkOSSID: "sess_1", Expiry: 1}, nil
	case "expired":
		return nil, services.ErrTokenExpired
	default:
		return nil, errors.New("bad token")
	}
}

func (f fakeAuth) RefreshSession(context.Context, string) (*value.SessionTokens, error) {
	if f.refreshErr != nil {
		return nil, f.refreshErr
	}
	return &value.SessionTokens{AccessToken: "fresh", RefreshToken: "fresh_refresh"}, nil
}

func newCtx(access, refresh string, auth services.AuthService) *AuthorizeCtx {
	ctx := &AuthorizeCtx{AccessToken: access, RefreshToken: refresh}
	ctx.Srv = &services.SrvHolder{AuthService: auth}
	return ctx
}

func assertUnauthorized(t *testing.T, err error) {
	t.Helper()
	var sc errs.StatusCoder
	if !errors.As(err, &sc) {
		t.Fatalf("error %v carries no status code", err)
	}
	if sc.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", sc.StatusCode())
	}
}

func TestResolveClaims_Valid(t *testing.T) {
	ctx := newCtx("valid", "", fakeAuth{})
	if err := (&resolveClaims{ctx: ctx}).Do(&core.DoCtx{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.claims == nil || ctx.claims.WorkOSUserID != "user_1" {
		t.Fatalf("claims not populated: %+v", ctx.claims)
	}
}

func TestResolveClaims_Invalid(t *testing.T) {
	ctx := newCtx("garbage", "", fakeAuth{})
	assertUnauthorized(t, (&resolveClaims{ctx: ctx}).Do(&core.DoCtx{}))
}

func TestRefresh_ExpiredAccessTokenIsRefreshed(t *testing.T) {
	ctx := newCtx("expired", "good_refresh", fakeAuth{})
	// resolveClaims leaves claims nil for an expired token.
	if err := (&resolveClaims{ctx: ctx}).Do(&core.DoCtx{}); err != nil {
		t.Fatalf("resolveClaims: %v", err)
	}
	if ctx.claims != nil {
		t.Fatal("expected nil claims after expired token")
	}
	if err := (&refreshIfExpired{ctx: ctx}).Do(&core.DoCtx{}); err != nil {
		t.Fatalf("refreshIfExpired: %v", err)
	}
	if ctx.NewTokens == nil || ctx.NewTokens.AccessToken != "fresh" {
		t.Fatalf("expected refreshed tokens, got %+v", ctx.NewTokens)
	}
	if ctx.claims == nil || ctx.claims.WorkOSUserID != "user_1" {
		t.Fatalf("expected claims from refreshed token, got %+v", ctx.claims)
	}
}

func TestRefresh_NoRefreshToken(t *testing.T) {
	ctx := newCtx("expired", "", fakeAuth{})
	_ = (&resolveClaims{ctx: ctx}).Do(&core.DoCtx{})
	assertUnauthorized(t, (&refreshIfExpired{ctx: ctx}).Do(&core.DoCtx{}))
}

func TestRefresh_RefreshFails(t *testing.T) {
	ctx := newCtx("expired", "good_refresh", fakeAuth{refreshErr: errors.New("revoked")})
	_ = (&resolveClaims{ctx: ctx}).Do(&core.DoCtx{})
	assertUnauthorized(t, (&refreshIfExpired{ctx: ctx}).Do(&core.DoCtx{}))
}
