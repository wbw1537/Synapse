package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// Database
	DBPath string `env:"SYNAPSE_DB_PATH" envDefault:"synapse.db"`

	// HTTP API
	HTTPPort string `env:"SYNAPSE_HTTP_PORT" envDefault:":8080"`

	// MQTT Broker
	MQTTPort string `env:"SYNAPSE_MQTT_PORT" envDefault:":1883"`
	WSPort   string `env:"SYNAPSE_WS_PORT" envDefault:":8083"` // WebSocket for UI

	// Security
	AuthToken string `env:"SYNAPSE_AUTH_TOKEN" envDefault:"synapse-secret"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	return cfg
}
