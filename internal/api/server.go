package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/models"
	"github.com/wbw1537/synapse/internal/service"
)

type Server struct {
	sessionsMu sync.Mutex
	sessions   map[string]time.Time
	cfg        *config.Config
	svcManager *service.Manager
	router     *chi.Mux
	staticFS   fs.FS
}

func NewServer(cfg *config.Config, svcManager *service.Manager, staticFS fs.FS) *Server {
	s := &Server{
		sessions:   make(map[string]time.Time),
		cfg:        cfg,
		svcManager: svcManager,
		router:     chi.NewRouter(),
		staticFS:   staticFS,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// API Routes
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Use(s.checkOrigin)
		r.Post("/session", s.login)
		r.Delete("/session", s.logout)
		r.Post("/discovery", s.registerService)
		r.Group(func(r chi.Router) {
			r.Use(s.requireOperator)
			r.Get("/session", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
			r.Get("/events", s.streamServices)
			r.Get("/services", s.listServices)
			r.Get("/services/{id}", s.getService)
			r.Post("/services/{id}/actions/{action_id}", s.executeAction)
		})
	})

	// Static Files (Frontend)
	if s.staticFS != nil {
		// We expect the FS to be rooted at web/dist
		distFS, err := fs.Sub(s.staticFS, "web/dist")
		if err != nil {
			log.Printf("Warning: Failed to locate web/dist in embedded FS: %v", err)
		} else {
			fileServer := http.FileServer(http.FS(distFS))
			s.router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If the path doesn't have an extension (not a file), serve index.html
				// This handles SPA routing
				if !strings.Contains(r.URL.Path, ".") && r.URL.Path != "/" {
					r.URL.Path = "/"
				}
				fileServer.ServeHTTP(w, r)
			}))
		}
	}
}

func (s *Server) Start() error {
	log.Printf("Starting HTTP API on %s", s.cfg.HTTPPort)
	server := &http.Server{Addr: s.cfg.HTTPPort, Handler: s.router, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	return server.ListenAndServe()
}

// Handlers

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.svcManager.List()
	if err != nil {
		http.Error(w, "Failed to list services", http.StatusInternalServerError)
		return
	}
	if services == nil {
		services = []models.Service{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	svc, err := s.svcManager.Get(id)
	if err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(svc)
}

func (s *Server) executeAction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actionID := chi.URLParam(r, "action_id")

	if err := s.svcManager.ExecuteAction(id, actionID); err != nil {
		log.Printf("ExecuteAction failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Action triggered"))
}

func (s *Server) registerService(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, models.MaxDiscoveryBytes))
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := s.svcManager.Upsert(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// streamServices sends only server-accepted state, including TTL transitions.
func (s *Server) streamServices(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	changes, cancel := s.svcManager.SubscribeChanges()
	defer cancel()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	send := func() error {
		if !s.authorized(r) {
			return fmt.Errorf("session expired")
		}
		services, err := s.svcManager.List()
		if err != nil {
			return err
		}
		if services == nil {
			services = []models.Service{}
		}
		data, err := json.Marshal(services)
		if err != nil {
			return err
		}
		// Bound writes so a slow/disconnected browser cannot retain a handler forever.
		_ = http.NewResponseController(w).SetWriteDeadline(time.Now().Add(10 * time.Second))
		if _, err := fmt.Fprintf(w, "event: services\ndata: %s\n\n", data); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}
	if err := send(); err != nil {
		return
	}
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-changes:
			if err := send(); err != nil {
				return
			}
		case <-ticker.C:
			if err := send(); err != nil {
				return
			}
		}
	}
}
