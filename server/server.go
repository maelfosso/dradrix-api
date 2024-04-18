// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"stockinos.com/api/broker"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

type Server struct {
	address  string
	database *storage.Database
	nats     *broker.Broker
	log      *zap.Logger
	mux      chi.Router
	server   *http.Server
}

type Options struct {
	Database *storage.Database
	Nats     *broker.Broker
	Host     string
	Log      *zap.Logger
	Port     int
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()

	return &Server{
		database: opts.Database, // createDatabase(opts.Log),
		nats:     createNats(opts.Log),
		address:  address,
		log:      opts.Log,
		mux:      mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func createDatabase(log *zap.Logger) *storage.Database {
	return storage.NewDatabase(storage.NewDatabaseOptions{
		URI:  utils.GetDefault("MONGODB_URI", "mongodb://localhost:27017/stockinos"),
		Name: utils.GetDefault("MONGDB_DBNAME", "stockinos"),
		Log:  log,
	})
}

func createNats(log *zap.Logger) *broker.Broker {
	return broker.NewBroker(broker.Options{
		Log:     log,
		Servers: "http://localhost:4222",
	})
}

// Start the server by setting up routes and listening for HTTP request on the given address
func (s *Server) Start() error {
	// if err := s.database.Connect(); err != nil {
	if s.database == nil {
		return fmt.Errorf("error with the database: not available")
	}

	// if err := s.nats.Connect(); err != nil {
	// 	return fmt.Errorf("error connecting to nats: %w", err)
	// }

	// if err := s.nats.Setup(); err != nil {
	// 	return fmt.Errorf("error setting up nats: %w", err)
	// }

	s.setupRoutes()

	// subscribers.NewMessageWoZSentSubscriber(*s.nats).Subscribe(*s.database)

	s.log.Info("Starting on", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout
func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
