package models

import (
	"time"
)

// Service represents the data model for a registered service
type Service struct {
	// Identity
	ID    string   `gorm:"primaryKey" json:"id"`
	Name  string   `gorm:"index" json:"name"`
	Group string   `gorm:"index" json:"group"`
	Tags  []string `gorm:"serializer:json" json:"tags"` // Serialized as JSON array in DB
	Icon  string   `json:"icon"`
	URL   string   `json:"url"`

	// State
	Status  string `gorm:"index" json:"status"` // online, warning, error, offline
	Message string `json:"message"`
	TTL     int    `json:"ttl"` // Seconds

	// Content
	Description  string `json:"description"`
	MarkdownDocs string `json:"markdown_docs"`

	// Dynamic Components
	Actions []Action `gorm:"serializer:json" json:"actions"`
	Widgets []Widget `gorm:"serializer:json" json:"widgets"`

	// Metadata
	LastSeen  time.Time `gorm:"index" json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Action represents an interactive button on the card
type Action struct {
	ID                  string `json:"id"`
	Label               string `json:"label"`
	Style               string `json:"style"` // primary, danger, default
	RequireConfirmation bool   `json:"require_confirmation"`
}

// Widget represents a UI element (graph, stat, etc.)
type Widget struct {
	Type     string                 `json:"type"`  // stat, progress, chart
	Label    string                 `json:"label"`
	Value    any                    `json:"value"`
	Unit     string                 `json:"unit"`
	Meta     map[string]interface{} `json:"meta"`
	Monitors []Monitor              `json:"monitors"`
}

// Monitor represents a server-side rule for alerting
type Monitor struct {
	Condition string `json:"condition"` // Expression: "value > 90"
	Severity  string `json:"severity"`  // "warning", "error"
	Message   string `json:"message"`   // "CPU High"
}

// ServicePayload is the exact structure expected from the MQTT discovery topic.
// It matches Service mostly, but includes the AuthToken.
type ServicePayload struct {
	APIVersion string `json:"api_version"`
	AuthToken  string `json:"auth_token"`
	Service           // Embed the Service struct fields
}
