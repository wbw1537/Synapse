package broker

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type Broker struct {
	Server *mqtt.Server
}

func New(axonToken, coreToken string) *Broker {
	// Create the new MQTT Server with default options
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "pk" || a.Key == "packet" {
			return slog.String(a.Key, "[redacted]")
		}
		return a
	}}))
	server := mqtt.New(&mqtt.Options{Logger: logger})
	server.Options.Capabilities.MaximumPacketSize = 1024*1024 + 4096
	server.Options.Capabilities.MaximumClients = 256

	_ = server.AddHook(&accessHook{axonToken: axonToken, coreToken: coreToken}, nil)

	return &Broker{Server: server}
}

func (b *Broker) Start(tcpPort, wsPort string) error {
	// TCP Listener
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: tcpPort,
	})
	if err := b.Server.AddListener(tcp); err != nil {
		return err
	}

	// WebSocket Listener
	ws := listeners.NewWebsocket(listeners.Config{
		ID:      "ws1",
		Address: wsPort,
	})
	if err := b.Server.AddListener(ws); err != nil {
		return err
	}

	// Start the broker in a separate goroutine so it doesn't block main immediately
	// However, the caller (main) usually handles the blocking via signals.
	// For now, we just start serving.
	go func() {
		log.Printf("Starting MQTT Broker on TCP %s and WS %s", tcpPort, wsPort)
		if err := b.Server.Serve(); err != nil {
			log.Fatalf("Error serving MQTT: %v", err)
		}
	}()

	return nil
}

func (b *Broker) Stop() {
	if err := b.Server.Close(); err != nil {
		log.Printf("Error closing MQTT broker: %v", err)
	}
}

// WaitForSignal is a helper to block until SIGINT/SIGTERM
func WaitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
