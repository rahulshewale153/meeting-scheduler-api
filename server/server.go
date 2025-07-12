package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahulshewale153/meeting-scheduler-api/configreader"
)

type server struct {
	httpServer *http.Server
	Config     *configreader.Config
}

func NewServer(config *configreader.Config) *server {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Connection.Port),
	}
	return &server{httpServer: httpServer, Config: config}
}

// service start with http endpoint
func (s *server) Start() {
	//setup http server
	r := mux.NewRouter()
	//basic health api
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	s.httpServer.Handler = r
	go func() {
		log.Println("Server starting on :8080")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	log.Println("server started...")
}

// service stop all http connection correctly, graceful shutdown occurred during running process
func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
