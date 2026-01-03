package models

import (
	"time"
)

// Service represents the data model for a registered service
// We define it here so it can be used by the DB, Broker, and API packages.
type Service struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"index" json:"name"`
	Status    string    `gorm:"index" json:"status"`
	LastSeen  time.Time `gorm:"index" json:"last_seen"`
	Payload   string    `gorm:"type:text" json:"-"` // Raw JSON payload, not exposed directly in JSON response usually
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
