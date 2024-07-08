package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"

	zlog "github.com/rs/zerolog/log"
)

type server struct {
	cfg *config

	// services
	todoSvc TodoService
}

func NewServer(todoSvc TodoService, options ...Option) (*server, error) {
	cfg := defaultConfig()
	for _, opt := range options {
		opt(cfg)
	}

	svr := &server{
		todoSvc: todoSvc,
		cfg:     cfg,
	}

	return svr, nil
}

func (s *server) Serve() {
	// First setup API router
	r := chi.NewRouter()
	s.setupMiddlewares(r)
	s.setupEndpoints(r)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlog.Fatal().Err(err).Msg("Error on srv.ListenAndServe()")
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	zlog.Print("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Printf("server shutdown returned an err: %v\n", err)
	}

	defer func() {
		// zlog.Print("Call s.db.Close()")
		// s.db.Close()
	}()

	zlog.Print("Server exiting")
}

func zerologCtx(env string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
			ctx := zerolog.New(output).With().Timestamp().Logger().WithContext(r.Context())

			if env == "prod" {
				ctx = zerolog.New(os.Stdout).WithContext(r.Context())
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (s *server) setupMiddlewares(r *chi.Mux) {
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Add 'zerolog' to request context
	r.Use(zerologCtx(s.cfg.Env))
}

func (s *server) setupEndpoints(r *chi.Mux) {
	// RESTy routes for "todos" resource
	r.Route("/api/v1/todos", func(r chi.Router) {
		r.Post("/", s.createTodo)       // POST /api/v1/todos
		r.Put("/{id}", s.updateTodo)    // PUT /api/v1/todos/1
		r.Get("/{id}", s.getTodo)       // GET /api/v1/todos/1
		r.Get("/", s.listTodos)         // GET /api/v1/todos?limit=10&offset=0
		r.Delete("/{id}", s.deleteTodo) // DELETE /api/v1/todos/1
	})
}
