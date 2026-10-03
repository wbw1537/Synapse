package api

import (
	"bufio"
	"context"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/db"
	"github.com/wbw1537/synapse/internal/service"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOperatorAccess(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	app := NewServer(&config.Config{AuthToken: "axon-secret", AdminToken: "operator-secret"}, service.NewManager(database, &config.Config{AuthToken: "axon-secret"}), nil)
	server := httptest.NewServer(app.router)
	defer server.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	request := func(method, path, body, origin string, want int) *http.Response {
		req, _ := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != want {
			res.Body.Close()
			t.Fatalf("%s %s: got %d want %d", method, path, res.StatusCode, want)
		}
		res.Body.Close()
		return res
	}
	for _, path := range []string{"/services", "/services/sample", "/events", "/session"} {
		request("GET", "/api/v1"+path, "", "", 401)
	}
	request("POST", "/api/v1/services/sample/actions/restart", "", "", 401)
	request("POST", "/api/v1/session", `{"token":"axon-secret"}`, "", 401)
	request("POST", "/api/v1/session", `{"token":"operator-secret"}`, "http://foreign.invalid", 403)
	res := request("POST", "/api/v1/session", `{"token":"operator-secret"}`, server.URL, 204)
	cookies := res.Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatal("session cookie policy missing")
	}
	request("GET", "/api/v1/services", "", "", 200)
	request("GET", "/api/v1/services", "", "http://foreign.invalid", 403)
	request("POST", "/api/v1/services/sample/actions/restart", "", "http://foreign.invalid", 403)
	request("DELETE", "/api/v1/session", "", "", 204)
	request("GET", "/api/v1/services", "", "", 401)
	request("POST", "/api/v1/session", `{"token":"operator-secret"}`, "", 204)
	app.sessionsMu.Lock()
	for key := range app.sessions {
		app.sessions[key] = time.Now().Add(-time.Second)
	}
	app.sessionsMu.Unlock()
	request("GET", "/api/v1/services", "", "", 401)
	bearer := &http.Client{Transport: authenticatedTransport{token: "operator-secret"}}
	res, err = bearer.Get(server.URL + "/api/v1/services")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatal("valid bearer denied")
	}
}

func TestExpiredSessionClosesStateStream(t *testing.T) {
	database, err := db.Connect(filepath.Join(t.TempDir(), "stream-auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.InitSchema(); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{AuthToken: "axon-secret", AdminToken: "operator-secret"}
	manager := service.NewManager(database, cfg)
	app := NewServer(cfg, manager, nil)
	server := httptest.NewServer(app.router)
	defer server.Close()
	login, err := http.Post(server.URL+"/api/v1/session", "application/json", strings.NewReader(`{"token":"operator-secret"}`))
	if err != nil {
		t.Fatal(err)
	}
	login.Body.Close()
	cookie := login.Cookies()[0]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+"/api/v1/events", nil)
	req.AddCookie(cookie)
	stream, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Body.Close()
	reader := bufio.NewReader(stream.Body)
	for i := 0; i < 3; i++ {
		if _, err := reader.ReadString('\n'); err != nil {
			t.Fatal(err)
		}
	}
	app.sessionsMu.Lock()
	app.sessions[cookie.Value] = time.Now().Add(-time.Second)
	app.sessionsMu.Unlock()
	if err := manager.Upsert([]byte(`{"api_version":"v1","auth_token":"axon-secret","id":"private","status":"online","ttl":30,"layout":{"type":"sections","root":[]},"components":{}}`)); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Fatalf("expired session received state: %s", data)
	}
}
