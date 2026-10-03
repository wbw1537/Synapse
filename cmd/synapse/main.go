package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

	// 5. Start MQTT Broker (Embedded)
	coreSecret := make([]byte, 32)
	if _, err := rand.Read(coreSecret); err != nil {
		log.Fatal("Cannot generate internal broker credential")
	}
	coreToken := hex.EncodeToString(coreSecret)
	mqttBroker := broker.New(cfg.AuthToken, coreToken)
	if err := mqttBroker.Start(cfg.MQTTPort, cfg.WSPort); err != nil {
		log.Fatalf("Failed to start MQTT broker: %v", err)
	}
	defer mqttBroker.Stop()

	// 6. Connect Internal MQTT Client (The "Core" Logic)
	// We wait a second to ensure the broker is fully up
	time.Sleep(1 * time.Second)

	opts := mqtt.NewClientOptions()
	host, port, err := net.SplitHostPort(cfg.MQTTPort)
	if err != nil {
		log.Fatalf("Invalid MQTT address: %v", err)
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	opts.AddBroker("tcp://" + net.JoinHostPort(host, port))
	opts.SetUsername("synapse-core")
	opts.SetPassword(coreToken)
	opts.SetClientID("synapse_core")
	opts.SetAutoReconnect(true)

	ready := make(chan struct{}, 1)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		topic := "synapse/v1/discovery/#"
		token := client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
			if err := svcManager.UpsertMQTT(msg.Topic(), msg.Payload()); err != nil {
				log.Printf("Discovery rejected: %v", err)
			}
		})
		if !token.WaitTimeout(5 * time.Second) {
			log.Print("Discovery subscription timed out")
			return
		}
		if token.Error() != nil {
			log.Printf("Discovery subscription failed: %v", token.Error())
			return
		}
		log.Printf("Listening for services on %s", topic)
		select {
		case ready <- struct{}{}:
		default:
		}
	})
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect internal MQTT client: %v", token.Error())
	}
	defer client.Disconnect(250)
	select {
	case <-ready:
	case <-time.After(6 * time.Second):
		log.Fatal("Initial discovery subscription unavailable")
	}

	// 6.5 Inject Publisher into Manager
	svcManager.SetPublisher(func(topic string, payload interface{}) error {
		// Serialize payload to JSON string/bytes
		importJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		token := client.Publish(topic, 1, false, importJSON)
		if !token.WaitTimeout(5 * time.Second) {
			return fmt.Errorf("command publication timed out")
		}
		return token.Error()
	})

	// Expose HTTP only after MQTT ingestion and command publication are ready.
	apiServer := api.NewServer(cfg, svcManager, synapse.UI)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("Failed to start HTTP API: %v", err)
		}
	}()

	log.Println("Synapse is running. Press Ctrl+C to stop.")

	// 8. Wait for shutdown signal
	broker.WaitForSignal()
	log.Println("Shutting down...")
}
