package broker

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mochi-mqtt/server/v2/listeners"
	"testing"
	"time"
)

func TestBrokerAccessBoundaries(t *testing.T) {
	b := New("axon-secret", "core-secret")
	listener := listeners.NewTCP(listeners.Config{ID: "test", Address: "127.0.0.1:0"})
	if err := b.Server.AddListener(listener); err != nil {
		t.Fatal(err)
	}
	if err := b.Server.Serve(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()
	wait := func(token paho.Token) error {
		if !token.WaitTimeout(2 * time.Second) {
			t.Fatal("MQTT timeout")
		}
		return token.Error()
	}
	connect := func(id, user, password string, allowed bool) paho.Client {
		c := paho.NewClient(paho.NewClientOptions().AddBroker("tcp://" + listener.Address()).SetClientID(id).SetUsername(user).SetPassword(password).SetAutoReconnect(false))
		err := wait(c.Connect())
		if allowed && err != nil {
			t.Fatal(err)
		}
		if !allowed && err == nil {
			t.Fatal("unauthorized broker connection succeeded")
		}
		t.Cleanup(func() { c.Disconnect(0) })
		return c
	}
	connect("anonymous", "", "", false)
	connect("wrong", "axon", "wrong", false)
	connect("synapse_core", "synapse-core", "axon-secret", false)
	axon := connect("sample", "axon", "axon-secret", true)
	core := connect("synapse_core", "synapse-core", "core-secret", true)
	discoveries := make(chan string, 4)
	commands := make(chan string, 4)
	if err := wait(core.Subscribe("synapse/v1/discovery/#", 1, func(_ paho.Client, msg paho.Message) { discoveries <- string(msg.Payload()) })); err != nil {
		t.Fatal(err)
	}
	if err := wait(axon.Subscribe("synapse/v1/command/sample", 1, func(_ paho.Client, msg paho.Message) { commands <- string(msg.Payload()) })); err != nil {
		t.Fatal(err)
	}
	for _, topic := range []string{"synapse/v1/discovery/#", "synapse/v1/command/#", "synapse/v1/command/other"} {
		token := axon.Subscribe(topic, 1, nil)
		wait(token)
		if token.(*paho.SubscribeToken).Result()[topic] != 0x80 {
			t.Fatalf("unauthorized subscription accepted: %s result=%v", topic, token.(*paho.SubscribeToken).Result())
		}
	}
	if err := wait(axon.Publish("synapse/v1/discovery/sample", 1, false, "permitted")); err != nil {
		t.Fatal(err)
	}
	select {
	case value := <-discoveries:
		if value != "permitted" {
			t.Fatal(value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("allowed discovery not delivered")
	}
	// PUBACK does not prove publication authorization; assert actual delivery.
	wait(axon.Publish("synapse/v1/discovery/other", 1, false, "forged"))
	wait(axon.Publish("synapse/v1/command/sample", 1, false, "forged-command"))
	select {
	case value := <-discoveries:
		t.Fatalf("forged discovery delivered: %s", value)
	case <-time.After(150 * time.Millisecond):
	}
	select {
	case value := <-commands:
		t.Fatalf("forged command delivered: %s", value)
	case <-time.After(150 * time.Millisecond):
	}
	// The broker disconnects MQTT v3 clients that publish forbidden topics.
	axon.Disconnect(0)
	axon = connect("sample", "axon", "axon-secret", true)
	if err := wait(axon.Subscribe("synapse/v1/command/sample", 1, func(_ paho.Client, msg paho.Message) { commands <- string(msg.Payload()) })); err != nil {
		t.Fatal(err)
	}
	if err := wait(core.Publish("synapse/v1/command/sample", 1, false, "permitted-command")); err != nil {
		t.Fatal(err)
	}
	select {
	case value := <-commands:
		if value != "permitted-command" {
			t.Fatal(value)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("allowed command not delivered")
	}
}
