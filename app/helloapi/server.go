package helloapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"

	"github.com/alexliesenfeld/health"

	zlog "github.com/rs/zerolog/log"
)

type server struct {
	opts *options

	// stores
	todoStore TodoStore
}

func NewServer(todoStore TodoStore, options ...Option) (*server, error) {
	opts := defaultConfig()
	for _, opt := range options {
		opt(opts)
	}

	svr := &server{
		todoStore: todoStore,
		opts:      opts,
	}

	return svr, nil
}

func (s *server) Serve(cleanup func()) {
	// First setup API router
	r := chi.NewRouter()
	s.setupMiddlewares(r)
	s.setupEndpoints(r)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.opts.Host, s.opts.Port),
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
		cleanup()
		zlog.Print("cleanup() function called")
	}()

	zlog.Print("Server exiting")
}

func zerologCtx(env string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
			ctx := zerolog.New(output).With().Timestamp().Logger().WithContext(r.Context())

			if strings.ToLower(env) == "prod" {
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
	r.Use(zerologCtx(s.opts.Env))
}

func (s *server) setupEndpoints(r *chi.Mux) {
	r.Route("/api/v1/todos", handleTodos(s))

	// health check
	checker := health.NewChecker()
	r.Handle("/health", health.NewHandler(checker))
}

func handleTodos(s *server) func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", makeAPIHandler(s.CreateTodo))
		r.Put("/{id}", makeAPIHandler(s.UpdateTodo))    // PUT /api/v1/todos/1
		r.Get("/{id}", makeAPIHandler(s.GetTodo))       // GET /api/v1/todos/1
		r.Get("/", makeAPIHandler(s.ListTodos))         // GET /api/v1/todos?limit=10&offset=0
		r.Delete("/{id}", makeAPIHandler(s.DeleteTodo)) // DELETE /api/v1/todos/1
	}
}
