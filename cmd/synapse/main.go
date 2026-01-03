package main

import (
	"log"

	"github.com/wbw1537/synapse/internal/broker"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
)

func main() {
	// 1. Load Config
	cfg := config.Load()
	log.Println("Synapse starting...")

	// 2. Initialize Database
	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// 3. Start MQTT Broker
	mqttBroker := broker.New()
	if err := mqttBroker.Start(cfg.MQTTPort, cfg.WSPort); err != nil {
		log.Fatalf("Failed to start MQTT broker: %v", err)
	}
	defer mqttBroker.Stop()
	log.Println("Synapse is running. Press Ctrl+C to stop.")

	// 4. Wait for shutdown signal
	broker.WaitForSignal()
	log.Println("Shutting down...")
}
