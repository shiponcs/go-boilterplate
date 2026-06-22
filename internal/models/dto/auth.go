package dto

// dto/ holds API request/response shapes only.

// SignupURLResponse is returned by GET /auth/signup — the AuthKit hosted URL
// the client should open to complete sign-up.
type SignupURLResponse struct {
	AuthorizationURL string `json:"authorization_url"`
}

// LogoutURLResponse is returned by GET /auth/logout — the WorkOS hosted URL
// the client should open to terminate the active AuthKit session.
type LogoutURLResponse struct {
	LogoutURL string `json:"logout_url"`
}

// UserResponse is the API shape of a signed-up user, returned by the callback.
type UserResponse struct {
	ID            uint   `json:"id"`
	WorkOSUserID  string `json:"workos_user_id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	EmailVerified bool   `json:"email_verified"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"created_at"`
}
