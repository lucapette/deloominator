package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api/graphql"
	"github.com/lucapette/deloominator/pkg/api/handlers"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
	"github.com/rs/cors"
	goji "goji.io"
	"goji.io/pat"
)

type Server struct {
	srv         *http.Server
	debug       bool
	dataSources db.DataSources
	port        string
	storage     *storage.Storage
}

type Option func(*Server)

func Port(port int) Option {
	return func(s *Server) {
		s.port = ":" + strconv.Itoa(port)
	}
}

func Debug(debug bool) Option {
	return func(s *Server) {
		s.debug = debug
	}
}

func DataSources(ds db.DataSources) Option {
	return func(s *Server) {
		s.dataSources = ds
	}
}

func Storage(storage *storage.Storage) Option {
	return func(s *Server) {
		s.storage = storage
	}
}

func NewServer(options []Option) *Server {
	s := &Server{}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// Start starts an api server according to the cfg configuration
func (s *Server) Start() {
	router := goji.NewMux()

	if s.debug {
		router.Use(handlers.Debug)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
	})

	if s.storage == nil {
		logrus.Warn("server running in read-only mode")
	}

	router.Use(handlers.Log)
	router.Use(c.Handler)

	router.HandleFunc(pat.Post("/graphql"), graphql.Handler(s.dataSources, s.storage))
	router.HandleFunc(pat.Get("/:name.:ext"), handlers.Static)
	router.HandleFunc(pat.Post("/export/:format"), handlers.Export(s.dataSources))
	router.HandleFunc(pat.Post("/query/evaluate"), handlers.QueryEvaluator)
	router.HandleFunc(pat.Get("/*"), handlers.UI)

	s.srv = &http.Server{Addr: s.port, Handler: router}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logrus.Printf("server closed")
			} else {
				logrus.Fatalf("cannot start server: %v", err)
			}
		}
	}()
}

// Stop shuts down the running API server
func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		logrus.Warnf("could not shutdown server: %v", err)
	}
}
