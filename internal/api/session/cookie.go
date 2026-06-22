// Package session reads and writes the sealed WorkOS session cookie. It is the
// single place that knows the cookie's name, attributes, and sealing, shared by
// the auth handler (which establishes the session) and the auth middleware
// (which validates and refreshes it).
package session

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/value"
	"github.com/example/go-svc-boilerplate/pkg/securecookie"
)

// cookieMaxAge is how long the session cookie persists in the browser. The
// access token inside is short-lived and refreshed on use; this bounds the
// overall session.
const cookieMaxAge = 14 * 24 * 60 * 60 // 14 days, seconds

// Set seals the tokens and writes the session cookie.
func Set(c *gin.Context, cfg *config.Config, tokens *value.SessionTokens) error {
	payload, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	sealed, err := securecookie.Seal([]byte(cfg.Session.CookiePassword), payload)
	if err != nil {
		return err
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cfg.Session.CookieName, sealed, cookieMaxAge, "/", cfg.Session.CookieDomain, cfg.Session.Secure, true)
	return nil
}

// Read returns the tokens from the session cookie, or an error if the cookie is
// absent or cannot be unsealed.
func Read(c *gin.Context, cfg *config.Config) (*value.SessionTokens, error) {
	raw, err := c.Cookie(cfg.Session.CookieName)
	if err != nil {
		return nil, err
	}
	payload, err := securecookie.Unseal([]byte(cfg.Session.CookiePassword), raw)
	if err != nil {
		return nil, err
	}
	var tokens value.SessionTokens
	if err := json.Unmarshal(payload, &tokens); err != nil {
		return nil, err
	}
	return &tokens, nil
}

// Clear deletes the session cookie.
func Clear(c *gin.Context, cfg *config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cfg.Session.CookieName, "", -1, "/", cfg.Session.CookieDomain, cfg.Session.Secure, true)
}
