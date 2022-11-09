package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/negroni"
	log "github.com/sirupsen/logrus"

	"admincheckapi/api/config"
	"admincheckapi/api/router"
	"admincheckapi/api/stat"
)

// Server stores all needed fields for an API server
type Server struct {
	server   *http.Server
	shutdown chan struct{}
}

// Runner is an interface for the API server
type Runner interface {
	Run()
}

//
// NewAPIServer initializes all seerver structures and starts shutdown listener
//
func NewAPIServer() (*Server, error) {
	log.Traceln("Begin: NewAPIServer")

	// basic negroni stuff init
	handler := negroni.New()
	router := router.NewRouter()
	handler.UseHandler(router)

	// set up of main server config structure
	server := &http.Server{
		Addr:           config.Setup.ServerIPAddress + ":" + config.Setup.ServerPort,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// server waits on it if interrupted
	shutdown := make(chan struct{})

	// Start shutdown handler as separate litener process
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint
		close(shutdown)
		log.Infoln("Server shutdown requested")
		os.Exit(0)
	}()

	log.Traceln("End: NewAPIServer")

	return &Server{server, shutdown}, nil
}

//
// Run starts the HTTP/HTTPS Sserver
//
func (s Server) Run() {
	log.Traceln("Begin: Run")

	log.Infoln("Starting HTTP server: " + s.server.Addr)
	stat.SetHealthy(stat.ServiceHealthy)
	stat.SetAlive(stat.ServiceAlive)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		stat.SetHealthy(stat.ServiceError)
		stat.SetAlive(stat.ServiceDead)
		log.Errorln("Error in ListenAndServe: " + err.Error())
	}

	// wait for shutdown signals
	<-s.shutdown

	log.Traceln("End: Run")
}
