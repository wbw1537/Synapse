package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wbw1537/synapse"
	"github.com/wbw1537/synapse/internal/api"
	"github.com/wbw1537/synapse/internal/broker"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/service"
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

	// 3. Initialize Service Manager
	svcManager := service.NewManager(database, cfg)
	// Start TTL Monitor (Run every 10 seconds)
	svcManager.StartTTLMonitor(10 * time.Second)

	// 4. Start HTTP API
	apiServer := api.NewServer(cfg, svcManager, synapse.UI)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("Failed to start HTTP API: %v", err)
		}
	}()

	// 5. Start MQTT Broker (Embedded)
	mqttBroker := broker.New()
	if err := mqttBroker.Start(cfg.MQTTPort, cfg.WSPort); err != nil {
		log.Fatalf("Failed to start MQTT broker: %v", err)
	}
	defer mqttBroker.Stop()

	// 6. Connect Internal MQTT Client (The "Core" Logic)
	// We wait a second to ensure the broker is fully up
	time.Sleep(1 * time.Second)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://localhost%s", cfg.MQTTPort))
	opts.SetClientID("synapse_core")
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect internal MQTT client: %v", token.Error())
	}
	defer client.Disconnect(250)

	// 7. Subscribe to Discovery Topic
	topic := "synapse/v1/discovery/#"
	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Log receipt (optional, verbose)
		// log.Printf("Received message on %s", msg.Topic())
		if err := svcManager.Upsert(msg.Payload()); err != nil {
			log.Printf("Error processing discovery payload: %v", err)
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("Listening for services on %s", topic)
	log.Println("Synapse is running. Press Ctrl+C to stop.")

	// 8. Wait for shutdown signal
	broker.WaitForSignal()
	log.Println("Shutting down...")
}
