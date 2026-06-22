package value

// SessionTokens is the WorkOS token pair that is sealed into the session
// cookie. The JSON tags are the on-the-wire shape inside the sealed blob.
type SessionTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SessionClaims are the fields read from a WorkOS access-token JWT: the WorkOS
// user id (sub), the WorkOS session id (sid, needed to build the logout URL),
// and the expiry (exp, unix seconds).
type SessionClaims struct {
	WorkOSUserID string
	WorkOSSID    string
	Expiry       int64
}
