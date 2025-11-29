package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/Godrik0/HackChange-Alpha/backend/docs"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/http/handlers"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	cfg           *config.Config
	router        *chi.Mux
	server        *http.Server
	logger        interfaces.Logger
	clientHandler *handlers.ClientHandler
}

func NewServer(cfg *config.Config, clientHandler *handlers.ClientHandler, logger interfaces.Logger) *Server {
	s := &Server{
		cfg:           cfg,
		logger:        logger,
		clientHandler: clientHandler,
	}

	s.setupRouter()

	s.server = &http.Server{
		Addr:         cfg.Server.GetServerAddr(),
		Handler:      s.router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupRouter() {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Recovery(s.logger))
	r.Use(middleware.LoggingMiddleware(s.logger))

	r.Use(chimiddleware.Timeout(60 * time.Second))

	r.Get("/health", s.healthCheckHandler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/clients", func(r chi.Router) {
			r.Get("/", s.clientHandler.ListClients)
			r.Get("/search", s.clientHandler.SearchClients)
			r.Post("/", s.clientHandler.CreateClient)
			r.Get("/{id}", s.clientHandler.GetClient)
			r.Put("/{id}", s.clientHandler.UpdateClient)
			r.Delete("/{id}", s.clientHandler.DeleteClient)
			r.Get("/{id}/scoring", s.clientHandler.CalculateScoring)
		})
	})

	s.router = r
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

// healthCheckHandler проверяет состояние сервиса
// @Summary      Проверка здоровья сервиса
// @Description  Возвращает статус работоспособности API
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}
