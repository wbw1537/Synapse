package api

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/models"
	"github.com/wbw1537/synapse/internal/service"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAuthoritativeEventSnapshots(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "events.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{AuthToken: "private-test-token", AdminToken: "operator-test-token"}
	manager := service.NewManager(database, cfg)
	client := &http.Client{Transport: authenticatedTransport{token: "operator-test-token"}}
	server := httptest.NewServer(NewServer(cfg, manager, nil).router)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	open := func() (*http.Response, *bufio.Reader) {
		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+"/api/v1/events", nil)
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.Header.Get("Content-Type") != "text/event-stream" {
			t.Fatal("missing SSE content type")
		}
		return res, bufio.NewReader(res.Body)
	}
	read := func(reader *bufio.Reader) []models.Service {
		event, err := reader.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		if event != "event: services\n" {
			t.Fatalf("unexpected event %q", event)
		}
		data, err := reader.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(data, cfg.AuthToken) || strings.Contains(data, "auth_token") {
			t.Fatal("stream leaked registration credential")
		}
		if _, err := reader.ReadString('\n'); err != nil {
			t.Fatal(err)
		}
		var services []models.Service
		if err := json.Unmarshal([]byte(strings.TrimPrefix(data, "data: ")), &services); err != nil {
			t.Fatal(err)
		}
		return services
	}
	res, reader := open()
	defer res.Body.Close()
	if got := read(reader); len(got) != 0 {
		t.Fatal("initial snapshot not empty")
	}
	payload := `{"api_version":"v1","auth_token":"private-test-token","id":"live","status":"online","ttl":30,"layout":{"type":"sections","root":[]},"components":{}}`
	posted, err := client.Post(server.URL+"/api/v1/discovery", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	posted.Body.Close()
	if posted.StatusCode != 200 {
		t.Fatal("HTTP discovery failed")
	}
	got := read(reader)
	if len(got) != 1 || got[0].Status != "online" {
		t.Fatal("HTTP discovery not streamed")
	}
	if err := manager.UpsertMQTT("synapse/v1/discovery/live", []byte(strings.Replace(payload, `"online"`, `"warning"`, 1))); err != nil {
		t.Fatal(err)
	}
	got = read(reader)
	if got[0].Status != "warning" {
		t.Fatal("MQTT accepted state not streamed")
	}
	res.Body.Close()
	// State can change while disconnected; reconnect must begin with current state.
	if err := manager.Upsert([]byte(strings.Replace(payload, `"online"`, `"error"`, 1))); err != nil {
		t.Fatal(err)
	}
	res2, reader2 := open()
	defer res2.Body.Close()
	got = read(reader2)
	if len(got) != 1 || got[0].Status != "error" {
		t.Fatal("reconnect did not reconcile missed state")
	}
}
