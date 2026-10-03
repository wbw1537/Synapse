package service

import (
	"encoding/json"
	"fmt"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func TestPersistedLogHistory(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "logs.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(database, &config.Config{AuthToken: "test"})
	send := func(value any, limit int) error {
		payload := map[string]any{"api_version": "v1", "auth_token": "test", "id": "logs", "status": "online", "ttl": 30,
			"layout":     map[string]any{"type": "sections", "root": []any{}},
			"components": map[string]any{"events": map[string]any{"id": "events", "type": "log_stream", "value": value, "max_items": limit}}}
		data, _ := json.Marshal(payload)
		return manager.Upsert(data)
	}
	cases := []struct {
		value any
		limit int
		want  []any
	}{
		{"first", 2, []any{"first"}}, {"second", 2, []any{"first", "second"}},
		{"third", 2, []any{"second", "third"}},
		{[]string{"a", "b", "c"}, 2, []any{"b", "c"}},
		{nil, 2, []any{"b", "c"}}, {"", 2, []any{"b", "c"}},
		{[]string{}, 2, []any{}}, {"after-clear", 2, []any{"after-clear"}},
	}
	for _, tc := range cases {
		if err := send(tc.value, tc.limit); err != nil {
			t.Fatal(err)
		}
		// Read through a fresh manager so the assertion checks SQLite, not memory.
		stored, err := NewManager(database, &config.Config{}).Get("logs")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stored.Components["events"].Value, tc.want) {
			t.Fatalf("value %v: got %#v want %#v", tc.value, stored.Components["events"].Value, tc.want)
		}
	}
	before, _ := manager.Get("logs")
	for _, bad := range []any{42, true, []any{"line", 42}, map[string]any{"line": "bad"}} {
		if err := send(bad, 2); err == nil {
			t.Fatalf("accepted invalid log %v", bad)
		}
	}
	if err := send("bad", -1); err == nil {
		t.Fatal("accepted negative retention")
	}
	after, _ := manager.Get("logs")
	if !reflect.DeepEqual(before, after) {
		t.Fatal("rejected log changed persisted state")
	}
	// First registration snapshots are bounded, even without an existing component.
	if err := database.Conn.Delete(before).Error; err != nil {
		t.Fatal(err)
	}
	if err := send([]string{"a", "b", "c"}, 2); err != nil {
		t.Fatal(err)
	}
	stored, _ := manager.Get("logs")
	if !reflect.DeepEqual(stored.Components["events"].Value, []any{"b", "c"}) {
		t.Fatal("first snapshot was not bounded")
	}
}

func TestStateChangesAndTTL(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(database, &config.Config{AuthToken: "test"})
	changes, cancel := manager.SubscribeChanges()
	defer cancel()
	payload := []byte(`{"api_version":"v1","auth_token":"test","id":"state","status":"online","ttl":1,"layout":{"type":"sections","root":[]},"components":{}}`)
	if err := manager.Upsert(payload); err != nil {
		t.Fatal(err)
	}
	select {
	case <-changes:
	default:
		t.Fatal("accepted HTTP state did not notify")
	}
	if err := manager.Upsert([]byte(`{"bad":true}`)); err == nil {
		t.Fatal("accepted invalid state")
	}
	select {
	case <-changes:
		t.Fatal("invalid state emitted change")
	default:
	}
	if err := database.Conn.Exec("UPDATE services SET last_seen = datetime('now', '-5 seconds')").Error; err != nil {
		t.Fatal(err)
	}
	manager.checkTTL()
	select {
	case <-changes:
	default:
		t.Fatal("TTL did not notify")
	}
	stored, _ := manager.Get("state")
	if stored.Status != "offline" {
		t.Fatal("TTL did not persist offline")
	}
	if err := manager.UpsertMQTT("synapse/v1/discovery/state", payload); err != nil {
		t.Fatal(err)
	}
	select {
	case <-changes:
	default:
		t.Fatal("MQTT recovery did not notify")
	}
	stored, _ = manager.Get("state")
	if stored.Status != "online" {
		t.Fatal("heartbeat did not restore online")
	}
	// A slow subscriber must not block mutations and catches up via full read.
	for i := 0; i < 20; i++ {
		if err := manager.Upsert(payload); err != nil {
			t.Fatal(err)
		}
	}
	select {
	case <-changes:
	default:
		t.Fatal("coalesced invalidation missing")
	}
}

func TestConcurrentLogEventsAreNotLost(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "concurrent.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(database, &config.Config{AuthToken: "test"})
	var workers sync.WaitGroup
	failures := make(chan error, 30)
	for i := 0; i < 30; i++ {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			payload := fmt.Sprintf(`{"api_version":"v1","auth_token":"test","id":"concurrent","status":"online","ttl":30,"layout":{"type":"sections","root":[]},"components":{"logs":{"id":"logs","type":"log_stream","value":"event-%d","max_items":30}}}`, i)
			if err := manager.Upsert([]byte(payload)); err != nil {
				failures <- err
			}
		}(i)
	}
	workers.Wait()
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	stored, err := manager.Get("concurrent")
	if err != nil {
		t.Fatal(err)
	}
	logs := stored.Components["logs"].Value.([]any)
	unique := map[any]bool{}
	for _, line := range logs {
		unique[line] = true
	}
	if len(logs) != 30 || len(unique) != 30 {
		t.Fatalf("concurrent events lost: %v", logs)
	}
}
