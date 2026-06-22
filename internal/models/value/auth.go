package value

// AuthnResult is the domain-level outcome of exchanging an AuthKit code. The
// services layer maps the WorkOS SDK response into this value object so core
// code never imports the SDK directly.
type AuthnResult struct {
	WorkOSUserID  string
	Email         string
	FirstName     string
	LastName      string
	EmailVerified bool

	// AccessToken/RefreshToken are the WorkOS session tokens. They are sealed
	// into the session cookie and never returned to the client directly.
	AccessToken  string
	RefreshToken string
}
