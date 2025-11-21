package server

import (
	"log"
	"net/http"
	"time"
	"tracker/pkg/api"
)

const webDir = "./web"

// Server structure
type Server struct {
	Logger     *log.Logger
	HttpServer *http.Server
}

// Http server constructor
func NewServer(logger *log.Logger, port string) *Server {

	// File server handler registration

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Registring all application handlers
	api.Init()

	// Server structure initialization

	server := &http.Server{
		Addr:         port,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger:     logger,
		HttpServer: server,
	}
}
