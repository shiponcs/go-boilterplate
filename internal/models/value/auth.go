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
}
