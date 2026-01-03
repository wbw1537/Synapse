package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/models"
	"gorm.io/gorm/clause"
)

type Manager struct {
	db     *db.Database
	config *config.Config
}

func NewManager(database *db.Database, cfg *config.Config) *Manager {
	return &Manager{
		db:     database,
		config: cfg,
	}
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
	// We extract the embedded Service struct
	svc := p.Service
	svc.LastSeen = time.Now()

	// 3. Upsert into DB
	// GORM Clause: OnConflict update all columns
	err := m.db.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&svc).Error

	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	log.Printf("Service registered: %s (%s)", svc.Name, svc.ID)
	return nil
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
	// Logic:
	// Find all services where status != 'offline'
	// AND (now - last_seen) > ttl
	// For SQLite: datetime(last_seen) < datetime('now', '-' || ttl || ' seconds')
	
	// Since GORM + SQLite with computed expiration is complex to write in pure Go,
	// we will run a raw SQL update query for efficiency.
	
	// NOTE: We assume 'ttl' column is in Seconds.
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
		// In the future: trigger "Service Lost" notification here
	}
}
