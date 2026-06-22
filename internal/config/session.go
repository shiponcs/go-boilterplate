package config

// SessionConfig controls the sealed session cookie. CookiePassword is the
// symmetric key (must be 32 bytes for AES-256) used to seal/unseal the cookie
// that carries the WorkOS tokens; it has no safe default and must be set via
// env. CookieName/CookieDomain/Secure shape the Set-Cookie attributes.
type SessionConfig struct {
	CookieName     string
	CookiePassword string
	CookieDomain   string
	Secure         bool
}
