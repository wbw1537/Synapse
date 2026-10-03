package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Database
	DBPath string `env:"SYNAPSE_DB_PATH" envDefault:"synapse.db"`

	// HTTP API
	HTTPPort string `env:"SYNAPSE_HTTP_PORT" envDefault:"127.0.0.1:8080"`

	// MQTT Broker
	MQTTPort string `env:"SYNAPSE_MQTT_PORT" envDefault:"127.0.0.1:1883"`
	WSPort   string `env:"SYNAPSE_WS_PORT" envDefault:"127.0.0.1:8083"` // Optional MQTT WebSocket clients

	// Security
	CookieSecure bool   `env:"SYNAPSE_COOKIE_SECURE" envDefault:"false"`
	AdminToken   string `env:"SYNAPSE_ADMIN_TOKEN"`
	AuthToken    string `env:"SYNAPSE_AUTH_TOKEN"`

	// Notification (SMTP)
	SMTPHost     string `env:"SYNAPSE_SMTP_HOST"`
	SMTPPort     string `env:"SYNAPSE_SMTP_PORT" envDefault:"587"`
	SMTPUser     string `env:"SYNAPSE_SMTP_USER"`
	SMTPPass     string `env:"SYNAPSE_SMTP_PASS"`
	SMTPFrom     string `env:"SYNAPSE_SMTP_FROM" envDefault:"synapse@localhost"`
	SMTPTo       string `env:"SYNAPSE_SMTP_TO"` // Comma separated list
	EnableAlerts bool   `env:"SYNAPSE_ENABLE_ALERTS" envDefault:"false"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	if cfg.AuthToken == "" || cfg.AdminToken == "" || cfg.AuthToken == cfg.AdminToken {
		log.Fatal("Set distinct nonempty SYNAPSE_AUTH_TOKEN and SYNAPSE_ADMIN_TOKEN")
	}
	return cfg
}
