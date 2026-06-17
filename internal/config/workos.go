package config

// WorkOSConfig holds the credentials for WorkOS AuthKit. ApiKey and ClientID
// come from the WorkOS dashboard; RedirectURI must match a callback URL
// registered there and points at this service's /auth/callback endpoint.
type WorkOSConfig struct {
	ApiKey      string
	ClientID    string
	RedirectURI string
}
