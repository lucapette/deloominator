package api

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api/graphql"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/rs/cors"
	goji "goji.io"
	"goji.io/pat"
)

type Server struct {
	debug       bool
	dataSources db.DataSources
	port        string
	storage     *db.Storage
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

func Storage(storage *db.Storage) Option {
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
		router.Use(debugHandler)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
	})

	if s.storage == nil {
		logrus.Warn("server running in read-only mode")
	}

	router.Use(logHandler)
	router.Use(c.Handler)

	router.HandleFunc(pat.Post("/graphql"), graphql.Handler(s.dataSources, s.storage))
	router.HandleFunc(pat.Get("/:name.:ext"), assetsHandler)
	router.HandleFunc(pat.Get("/*"), uiHandler)

	go func() {
		if err := http.ListenAndServe(s.port, router); err != nil {
			logrus.Fatalf("cannot start server: %v", err)
		}
	}()
}
