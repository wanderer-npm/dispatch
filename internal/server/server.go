package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wanderer-npm/dispatch/internal/config"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/webhook", &webhookHandler{cfg: s.cfg})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := fmt.Sprintf(":%s", s.cfg.Server.Port)
	log.Printf("[dispatch] listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}
