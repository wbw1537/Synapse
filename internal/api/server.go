package api

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/wbw1537/synapse/internal/config"
	"github.com/wbw1537/synapse/internal/service"
)

type Server struct {
	cfg        *config.Config
	svcManager *service.Manager
	router     *chi.Mux
	staticFS   fs.FS
}

func NewServer(cfg *config.Config, svcManager *service.Manager, staticFS fs.FS) *Server {
	s := &Server{
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
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all for MVP
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	// API Routes
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/services", s.listServices)
		r.Get("/services/{id}", s.getService)
		r.Post("/services/{id}/actions/{action_id}", s.executeAction)
		r.Post("/discovery", s.registerService)
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
	return http.ListenAndServe(s.cfg.HTTPPort, s.router)
}

// Handlers

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.svcManager.List()
	if err != nil {
		http.Error(w, "Failed to list services", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(services)
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	svc, err := s.svcManager.Get(id)
	if err != nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
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
	body, err := io.ReadAll(r.Body)
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
