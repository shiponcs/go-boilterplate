package entity

import "time"

// User is the local mirror of a WorkOS-managed user. WorkOS is the source of
// truth for identity; this row links a WorkOS user to local domain data via
// WorkOSUserID. entity/ holds DB-mapped types only — no business logic.
type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	WorkOSUserID  string    `gorm:"column:workos_user_id;uniqueIndex" json:"workos_user_id"`
	Email         string    `gorm:"column:email;index" json:"email"`
	FirstName     string    `gorm:"column:first_name" json:"first_name"`
	LastName      string    `gorm:"column:last_name" json:"last_name"`
	EmailVerified bool      `gorm:"column:email_verified" json:"email_verified"`
	Status        string    `gorm:"column:status" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
