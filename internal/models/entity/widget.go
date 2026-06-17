package entity

import "time"

// Widget is an example GORM row. entity/ holds DB-mapped types only — no
// business logic. Replace Widget with your own domain entities.
type Widget struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Status    string    `gorm:"column:status" json:"status"`
	Price     float64   `gorm:"column:price" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Widget) TableName() string {
	return "widgets"
}
