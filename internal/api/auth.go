package api

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const sessionCookie = "synapse_session"
const sessionDuration = 12 * time.Hour

func secretMatches(a, b string) bool {
	return a != "" && b != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Origin checking complements SameSite cookies, including login and discovery.
func (s *Server) checkOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			u, err := url.Parse(origin)
			expectedScheme := "http"
			if r.TLS != nil || s.cfg.CookieSecure {
				expectedScheme = "https"
			}
			if err != nil || u.Host != r.Host || u.Scheme != expectedScheme {
				http.Error(w, "origin forbidden", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func (s *Server) authorized(r *http.Request) bool {
	if bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "); bearer != r.Header.Get("Authorization") && secretMatches(bearer, s.cfg.AdminToken) {
		return true
	}
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return false
	}
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	expires, ok := s.sessions[cookie.Value]
	if !ok || !time.Now().Before(expires) {
		delete(s.sessions, cookie.Value)
		return false
	}
	return true
}
func (s *Server) requireOperator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.authorized(r) {
			http.Error(w, "operator authentication required", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil || !secretMatches(input.Token, s.cfg.AdminToken) {
		http.Error(w, "invalid operator credential", http.StatusUnauthorized)
		return
	}
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		http.Error(w, "session unavailable", 500)
		return
	}
	token := hex.EncodeToString(random)
	expires := time.Now().Add(sessionDuration)
	s.sessionsMu.Lock()
	for key, expiry := range s.sessions {
		if !time.Now().Before(expiry) {
			delete(s.sessions, key)
		}
	}
	s.sessions[token] = expires
	s.sessionsMu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: token, Path: "/api/v1", HttpOnly: true, Secure: s.cfg.CookieSecure || r.TLS != nil, SameSite: http.SameSiteStrictMode, Expires: expires, MaxAge: int(sessionDuration.Seconds())})
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusNoContent)
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		s.sessionsMu.Lock()
		delete(s.sessions, cookie.Value)
		s.sessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Path: "/api/v1", HttpOnly: true, Secure: s.cfg.CookieSecure || r.TLS != nil, SameSite: http.SameSiteStrictMode, MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
}
