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
	Tags  []string `gorm:"serializer:json" json:"tags"`
	Icon  string   `json:"icon"`
	URL   string   `json:"url"`

	// State
	Status  string `gorm:"index" json:"status"` // online, warning, error, offline
	Message string `json:"message"`
	TTL     int    `json:"ttl"` // Seconds

	// Content
	Description  string `json:"description"`
	MarkdownDocs string `json:"markdown_docs"`

	// Layout & Components (Protocol v2)
	APIVersion string                  `json:"api_version"`
	Layout     LayoutSchema            `gorm:"serializer:json" json:"layout"`
	Components map[string]Component    `gorm:"serializer:json" json:"components"`

	// Metadata
	LastSeen  time.Time `gorm:"index" json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LayoutSchema defines the visual hierarchy
type LayoutSchema struct {
	Type string          `json:"type"` // e.g. "sections"
	Root []LayoutSection `json:"root"`
}

// LayoutSection defines a group of components
type LayoutSection struct {
	Type     string   `json:"type"` // "section"
	Title    string   `json:"title"`
	Children []string `json:"children"` // IDs referencing Components
}

// Component represents a UI widget/element
// Formerly "Widget", now stored in a map. Properties are flat for simplicity.
type Component struct {
	ID    string `json:"id"`
	Type  string `json:"type"` // stat, status_indicator, gauge, log_stream, action_group, link
	Label string `json:"label"`
	Value any    `json:"value"`
	Unit  string `json:"unit"`

	// Specific fields (Union of all widget properties)
	Copyable   bool                   `json:"copyable,omitempty"`
	Mapping    map[string]StatusState `json:"mapping,omitempty"`
	Min        float64                `json:"min,omitempty"`
	Max        float64                `json:"max,omitempty"`
	Thresholds map[string]string      `json:"thresholds,omitempty"`
	MaxItems   int                    `json:"max_items,omitempty"`
	Items      []ActionGroupItem      `json:"items,omitempty"`
	ActionID   string                 `json:"action_id,omitempty"` // For standalone buttons (if any)
	URI        string                 `json:"uri,omitempty"`
	Text       string                 `json:"text,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Animate    bool                   `json:"animate,omitempty"`
	Style      string                 `json:"style,omitempty"`
	Confirm    bool                   `json:"confirm,omitempty"`

	Monitors []Monitor `json:"monitors"`
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
	Condition string `json:"condition"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
}

// ServicePayload matches the MQTT discovery JSON
type ServicePayload struct {
	APIVersion string `json:"api_version"`
	AuthToken  string `json:"auth_token"`
	Service           // Embed Service fields
}
