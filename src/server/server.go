package server

import (
	"github.com/phuwn/crawlie/src/auth"
	"github.com/phuwn/crawlie/src/config"
	"github.com/phuwn/crawlie/src/database"
	"github.com/phuwn/crawlie/src/service"
	"github.com/phuwn/crawlie/src/store"
)

var srv *Server

// Server - server structure
type Server struct {
	auth     *auth.Authenticator
	database database.Database
	store    *store.Store
	service  *service.Service
}

// Init - create new server
func Init(cfg *config.Config) error {
	var (
		db  database.Database
		err error
	)

	switch cfg.Type {
	case "postgres":
		db, err = database.NewPostgresConn(cfg.Database)
	}
	if err != nil {
		return err
	}

	auth, err := auth.NewAuthenticator(cfg.Authenticator)
	if err != nil {
		return err
	}

	srv = &Server{
		auth:     auth,
		database: db,
		store:    store.New(),
		service:  service.New(cfg.Service),
	}
	return nil
}

// GetServer - get server configuration settings
func Get() *Server {
	return srv
}

// Auth - get Authenticator
func (s *Server) Auth() *auth.Authenticator {
	return s.auth
}

// Database - get Database
func (s *Server) DB() database.Database {
	return s.database
}

// Store - get store
func (s *Server) Store() *store.Store {
	return s.store
}

// Service - get service
func (s *Server) Service() *service.Service {
	return s.service
}
