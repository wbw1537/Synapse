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
	ID      string `json:"id"`
	Type    string `json:"type"` // stat, status_indicator, gauge, log_stream, action_group, link
	Label   string `json:"label"`
	Value   any    `json:"value"`
	Unit    string `json:"unit"`
	Visible *bool  `json:"visible,omitempty"`

	// Specific fields for various widgets
	Copyable   bool                   `json:"copyable,omitempty"`
	Mapping    map[string]StatusState `json:"mapping,omitempty"`
	Min        float64                `json:"min,omitempty"`
	Max        float64                `json:"max,omitempty"`
	Thresholds map[string]string      `json:"thresholds,omitempty"`
	MaxItems   int                    `json:"max_items,omitempty"`
	Items      []ActionGroupItem      `json:"items,omitempty"`
	ActionID   string                 `json:"action_id,omitempty"`
	URI        string                 `json:"uri,omitempty"`
	Text       string                 `json:"text,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Animate    bool                   `json:"animate,omitempty"`
	Style      string                 `json:"style,omitempty"`
	Confirm    bool                   `json:"confirm,omitempty"`

	Meta     map[string]interface{} `json:"meta,omitempty"`
	Monitors []Monitor              `json:"monitors"`
}

// StatusState defines the visual style for a status_indicator state
type StatusState struct {
	Text    string `json:"text"`
	Color   string `json:"color"`
	Icon    string `json:"icon"`
	Animate bool   `json:"animate,omitempty"`
}

// ActionGroupItem defines a button in an action_group
type ActionGroupItem struct {
	ActionID string `json:"action_id"`
	Label    string `json:"label"`
	Style    string `json:"style"`
	Confirm  bool   `json:"confirm"`
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
