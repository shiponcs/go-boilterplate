package cnst

// Domain constants live here so magic strings/values are defined once.

const (
	WidgetStatusActive   = "active"
	WidgetStatusInactive = "inactive"
)

const (
	UserStatusActive = "active"
)

// Language codes used by the localization layer.
const (
	LangEN = "en"
)

// Gin context keys for values stashed by middleware.
const (
	CtxUserKey = "current_user"
)
