package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/wbw1537/synapse/internal/broker"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/models"
	"github.com/wbw1537/synapse/internal/service"
)

func TestDiscoveryTransports(t *testing.T) {
	fixture, err := os.ReadFile("../../testdata/discovery-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	database, err := db.Connect(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{AuthToken: "fixture-token", AdminToken: "operator-test-token"}
	manager := service.NewManager(database, cfg)
	client := &http.Client{Transport: authenticatedTransport{token: "operator-test-token"}}
	server := httptest.NewServer(NewServer(cfg, manager, nil).router)
	defer server.Close()

	b := broker.New("fixture-token", "core-test-token")
	tcp := listeners.NewTCP(listeners.Config{ID: "test", Address: "127.0.0.1:0"})
	if err := b.Server.AddListener(tcp); err != nil {
		t.Fatal(err)
	}
	if err := b.Server.Serve(); err != nil {
		t.Fatal(err)
	}
	defer b.Stop()
	connect := func(id string) mqtt.Client {
		opts := mqtt.NewClientOptions().AddBroker("tcp://" + tcp.Address()).SetClientID(id)
		if id == "synapse_core" {
			opts.SetUsername("synapse-core").SetPassword("core-test-token")
		} else {
			opts.SetUsername("axon").SetPassword("fixture-token")
		}
		c := mqtt.NewClient(opts)
		waitToken(t, c.Connect())
		t.Cleanup(func() { c.Disconnect(0) })
		return c
	}
	core, axon := connect("synapse_core"), connect("protocol-demo")
	manager.SetPublisher(func(topic string, payload interface{}) error {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		token := core.Publish(topic, 1, false, data)
		if !token.WaitTimeout(3 * time.Second) {
			return &responseError{"command publish timeout"}
		}
		return token.Error()
	})
	commands := make(chan []byte, 1)
	waitToken(t, axon.Subscribe("synapse/v1/command/protocol-demo", 1, func(_ mqtt.Client, msg mqtt.Message) {
		commands <- append([]byte(nil), msg.Payload()...)
	}))
	results := make(chan error, 1)
	waitToken(t, core.Subscribe("synapse/v1/discovery/#", 1, func(_ mqtt.Client, msg mqtt.Message) {
		results <- manager.UpsertMQTT(msg.Topic(), msg.Payload())
	}))
	send := func(transport, topic string, payload []byte) error {
		if transport == "mqtt" {
			waitToken(t, axon.Publish(topic, 1, false, payload))
			select {
			case err := <-results:
				return err
			case <-time.After(3 * time.Second):
				t.Fatal("discovery callback not received")
			}
		}
		res, err := client.Post(server.URL+"/api/v1/discovery", "application/json", bytes.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK {
			return &responseError{string(body)}
		}
		return nil
	}
	var expected models.ServicePayload
	if err := json.Unmarshal(fixture, &expected); err != nil {
		t.Fatal(err)
	}
	for _, transport := range []string{"http", "mqtt"} {
		t.Run(transport, func(t *testing.T) {
			if err := send(transport, "synapse/v1/discovery/protocol-demo", fixture); err != nil {
				t.Fatal(err)
			}
			res, err := client.Get(server.URL + "/api/v1/services/protocol-demo")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			var stored models.Service
			if err := json.NewDecoder(res.Body).Decode(&stored); err != nil {
				t.Fatal(err)
			}
			if stored.APIVersion != "v1" || !reflect.DeepEqual(stored.Layout, expected.Layout) || !reflect.DeepEqual(stored.Components, expected.Components) {
				t.Fatalf("stored/UI structure differs: %+v", stored)
			}
			// The stored action_id drives the actual MQTT command envelope.
			res, err = client.Post(server.URL+"/api/v1/services/protocol-demo/actions/restart", "application/json", nil)
			if err != nil {
				t.Fatal(err)
			}
			res.Body.Close()
			if res.StatusCode != http.StatusOK {
				t.Fatalf("action request status %d", res.StatusCode)
			}
			select {
			case data := <-commands:
				var command map[string]string
				if err := json.Unmarshal(data, &command); err != nil || command["action_id"] != "restart" || command["issued_by"] != "synapse-ui" {
					t.Fatalf("invalid command envelope: %s", data)
				}
				if _, err := time.Parse(time.RFC3339, command["timestamp"]); err != nil {
					t.Fatal(err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("declared action command not received")
			}
			before, err := manager.Get(expected.ID)
			if err != nil {
				t.Fatal(err)
			}
			cases := []struct {
				name, want string
				edit       func(map[string]any)
			}{
				{"token", "invalid auth_token", func(p map[string]any) { p["auth_token"] = "wrong" }},
				{"missing token", "invalid auth_token", func(p map[string]any) { delete(p, "auth_token") }},
				{"missing version", "api_version", func(p map[string]any) { delete(p, "api_version") }},
				{"unsupported version", "api_version", func(p map[string]any) { p["api_version"] = "v2" }},
				{"numeric version", "invalid discovery", func(p map[string]any) { p["api_version"] = 1 }},
				{"ghost", "unknown component", func(p map[string]any) {
					p["layout"].(map[string]any)["root"].([]any)[0].(map[string]any)["children"] = []string{"ghost"}
				}},
				{"component id", "identical id", func(p map[string]any) { p["components"].(map[string]any)["cpu"].(map[string]any)["id"] = "other" }},
				{"missing component id", "identical id", func(p map[string]any) { delete(p["components"].(map[string]any)["cpu"].(map[string]any), "id") }},
				{"unknown type", "unsupported type", func(p map[string]any) { p["components"].(map[string]any)["cpu"].(map[string]any)["type"] = "button" }},
				{"missing layout", "layout requires", func(p map[string]any) { delete(p, "layout") }},
				{"null components", "components must", func(p map[string]any) { p["components"] = nil }},
				{"null children", "children array", func(p map[string]any) {
					p["layout"].(map[string]any)["root"].([]any)[0].(map[string]any)["children"] = nil
				}},
				{"layout type", "layout requires", func(p map[string]any) { p["layout"].(map[string]any)["type"] = "grid" }},
				{"unsafe id", "service id", func(p map[string]any) { p["id"] = "bad/id" }},
				{"missing action", "action_id", func(p map[string]any) {
					delete(p["components"].(map[string]any)["controls"].(map[string]any)["items"].([]any)[0].(map[string]any), "action_id")
				}},
				{"TOML action id", "unknown field", func(p map[string]any) {
					item := p["components"].(map[string]any)["controls"].(map[string]any)["items"].([]any)[0].(map[string]any)
					delete(item, "action_id")
					item["id"] = "restart"
				}},
				{"duplicate action", "unique", func(p map[string]any) {
					c := p["components"].(map[string]any)["controls"].(map[string]any)
					c["items"] = append(c["items"].([]any), c["items"].([]any)[0])
				}},
				{"standalone action", "action_group.items", func(p map[string]any) {
					p["components"].(map[string]any)["cpu"].(map[string]any)["action_id"] = "hidden"
				}},
				{"widgets", "legacy widgets", func(p map[string]any) { p["widgets"] = []any{} }},
				{"legacy actions", "legacy actions", func(p map[string]any) { p["actions"] = []any{} }},
				{"meta", "unknown field", func(p map[string]any) { p["meta"] = map[string]any{"id": "nested"} }},
				{"props", "unknown field", func(p map[string]any) {
					p["components"].(map[string]any)["cpu"].(map[string]any)["props"] = map[string]any{"max": 100}
				}},
				{"bad ttl", "ttl", func(p map[string]any) { p["ttl"] = 0 }},
				{"bad status", "status", func(p map[string]any) { p["status"] = "healthy" }},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					var p map[string]any
					json.Unmarshal(fixture, &p)
					tc.edit(p)
					payload, _ := json.Marshal(p)
					err := send(transport, "synapse/v1/discovery/protocol-demo", payload)
					if err == nil || !strings.Contains(err.Error(), tc.want) {
						t.Fatalf("want %q error, got %v", tc.want, err)
					}
					after, err := manager.Get(expected.ID)
					if err != nil || !reflect.DeepEqual(before, after) {
						t.Fatalf("rejected payload changed stored state: %v", err)
					}
				})
			}
			for _, payload := range []string{"null", "[]", "{", string(fixture) + " {}"} {
				if err := send(transport, "synapse/v1/discovery/protocol-demo", []byte(payload)); err == nil {
					t.Fatal("accepted malformed JSON")
				}
			}
		})
	}
	for _, topic := range []string{"synapse/v1/discovery/other", "synapse/v1/discovery/", "synapse/v1/discovery/protocol-demo/extra"} {
		if err := manager.UpsertMQTT(topic, fixture); err == nil {
			t.Fatalf("accepted topic %s", topic)
		}
	}
	// Empty snapshots are valid; every accepted update replaces capability definitions.
	empty := []byte(`{"api_version":"v1","auth_token":"fixture-token","id":"protocol-demo","status":"online","ttl":30,"layout":{"type":"sections","root":[]},"components":{}}`)
	for _, transport := range []string{"http", "mqtt"} {
		if err := send(transport, "synapse/v1/discovery/protocol-demo", empty); err != nil {
			t.Fatal(err)
		}
		stored, err := manager.Get("protocol-demo")
		if err != nil || len(stored.Components) != 0 || len(stored.Layout.Root) != 0 {
			t.Fatalf("snapshot failed to remove capabilities: %v", err)
		}
	}
}

type responseError struct{ message string }

func (e *responseError) Error() string { return e.message }
func waitToken(t *testing.T, token mqtt.Token) {
	t.Helper()
	if !token.WaitTimeout(3 * time.Second) {
		t.Fatal("MQTT timeout")
	}
	if err := token.Error(); err != nil {
		t.Fatal(err)
	}
}

type authenticatedTransport struct{ token string }

func (a authenticatedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+a.token)
	return http.DefaultTransport.RoundTrip(r)
}
