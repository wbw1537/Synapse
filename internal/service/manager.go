package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/evaluator"
	"github.com/wbw1537/synapse/internal/models"
	"github.com/wbw1537/synapse/internal/notification"
	"gorm.io/gorm/clause"
)

type Manager struct {
	db           *db.Database
	config       *config.Config
	alertManager *notification.AlertManager
	publishFunc  func(topic string, payload interface{}) error
}

func NewManager(database *db.Database, cfg *config.Config) *Manager {
	sender := notification.NewSMTPSender(cfg)
	return &Manager{
		db:           database,
		config:       cfg,
		alertManager: notification.NewAlertManager(sender),
	}
}

// SetPublisher sets the MQTT publish function
func (m *Manager) SetPublisher(fn func(topic string, payload interface{}) error) {
	m.publishFunc = fn
}

// ExecuteAction triggers a command to the remote axon
func (m *Manager) ExecuteAction(serviceID, actionID string) error {
	svc, err := m.Get(serviceID)
	if err != nil {
		return err
	}
	
	// 1. Validate if action exists in the service definition
	// We search through all components for action groups containing this ID
	found := false
	
	for _, comp := range svc.Components {
		if comp.Type == "action_group" {
			for _, item := range comp.Items {
				if item.ActionID == actionID {
					found = true
					break
				}
			}
		}
		// Also check if the component itself IS the action (if defined that way in future)
		if comp.ActionID == actionID {
			found = true
		}
		
		if found { break }
	}
	
	if !found {
		return fmt.Errorf("action '%s' not found for service '%s'", actionID, serviceID)
	}
	
	// 2. Publish Command
	if m.publishFunc == nil {
		return fmt.Errorf("mqtt publisher not configured")
	}
	
	topic := fmt.Sprintf("synapse/v1/command/%s", serviceID)
	payload := map[string]string{
		"action_id": actionID,
		"issued_by": "synapse-ui",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	return m.publishFunc(topic, payload)
}

// Upsert handles the registration/update logic
func (m *Manager) Upsert(payload []byte) error {
	var p models.ServicePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	// 1. Validation
	if p.ID == "" {
		return fmt.Errorf("service id is required")
	}
	if p.AuthToken != m.config.AuthToken {
		return fmt.Errorf("invalid auth_token")
	}

	// 2. Prepare Model
	svc := p.Service
	svc.LastSeen = time.Now()

	// 2.5 Merge with existing state (for log_stream, etc.)
	if existing, err := m.Get(svc.ID); err == nil && existing != nil {
		m.mergeComponents(existing, &svc)
	}

	// 3. Upsert into DB
	err := m.db.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&svc).Error

	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	// 4. Check Monitors
	m.evaluateMonitors(&svc)

	log.Printf("Service registered/updated: %s (%s)", svc.Name, svc.ID)
	return nil
}

// mergeComponents preserves state from existing components into the new update
func (m *Manager) mergeComponents(existing, incoming *models.Service) {
	if existing.Components == nil {
		return
	}

	for id, newComp := range incoming.Components {
		// 1. Handle Log Stream
		if newComp.Type == "log_stream" {
			if oldComp, ok := existing.Components[id]; ok {
				// Initialize or cast existing logs
				var logs []interface{}
				
				// Handle different potential types from JSON unmarshalling
				switch v := oldComp.Value.(type) {
				case []interface{}:
					logs = v
				case []string:
					for _, s := range v {
						logs = append(logs, s)
					}
				}

				// Append new value if it's a string
				if newVal, ok := newComp.Value.(string); ok {
					logs = append(logs, newVal)
				}

				// Enforce MaxItems
				maxItems := 10 // Default
				if newComp.MaxItems > 0 {
					maxItems = newComp.MaxItems
				}

				if len(logs) > maxItems {
					logs = logs[len(logs)-maxItems:]
				}

				newComp.Value = logs
				// Must write back to map because 'newComp' is a copy/loop variable value in Go maps? 
				// Actually range over map gives value copy. So we need to reassign.
				incoming.Components[id] = newComp
			} else {
				// First time seeing this log stream, wrap the single string in a list
				if val, ok := newComp.Value.(string); ok {
					newComp.Value = []string{val}
					incoming.Components[id] = newComp
				}
			}
		}
	}
}

func (m *Manager) evaluateMonitors(svc *models.Service) {
	for compID, comp := range svc.Components {
		for mIdx, monitor := range comp.Monitors {
			triggered, err := evaluator.Evaluate(monitor.Condition, comp.Value)
			if err != nil {
				log.Printf("Monitor evaluation error (svc=%s, comp=%s): %v", svc.ID, compID, err)
				continue
			}

			// Unique key for state tracking
			key := fmt.Sprintf("%s:%s:m%d", svc.ID, compID, mIdx)
			m.alertManager.CheckAndAlert(key, triggered, monitor.Severity, monitor.Message, svc.Name)
		}
	}
}

// StartTTLMonitor checks for expired services
func (m *Manager) StartTTLMonitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			m.checkTTL()
		}
	}()
}

func (m *Manager) checkTTL() {
	query := `
		UPDATE services 
		SET status = 'offline', updated_at = ?
		WHERE status != 'offline' 
		AND datetime(last_seen) < datetime(?, '-' || ttl || ' seconds')
	`
	
	now := time.Now()
	result := m.db.Conn.Exec(query, now, now)
	
	if result.Error != nil {
		log.Printf("Error checking TTL: %v", result.Error)
		return
	}
	
	if result.RowsAffected > 0 {
		log.Printf("Marked %d services as offline", result.RowsAffected)
	}
}

// List returns all services
func (m *Manager) List() ([]models.Service, error) {
	var services []models.Service
	result := m.db.Conn.Find(&services)
	return services, result.Error
}

// Get returns a single service by ID
func (m *Manager) Get(id string) (*models.Service, error) {
	var svc models.Service
	result := m.db.Conn.First(&svc, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &svc, nil
}
