package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rahulshewale153/meeting-scheduler-api/configreader"
	"github.com/rahulshewale153/meeting-scheduler-api/handler"
	"github.com/rahulshewale153/meeting-scheduler-api/repository"
	"github.com/rahulshewale153/meeting-scheduler-api/service"
)

type server struct {
	httpServer *http.Server
	config     *configreader.Config
	mysqlDB    *sql.DB
}

func NewServer(config *configreader.Config) *server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Connection.Port),
		ReadTimeout:  time.Duration(config.Connection.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Connection.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Connection.IdleTimeout) * time.Second,
	}

	sqlConn, err := setupMysqlDBConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	return &server{httpServer: httpServer, config: config, mysqlDB: sqlConn}
}

func setupMysqlDBConnection(config *configreader.Config) (*sql.DB, error) {
	mysqlConfig := config.MySQL
	fmt.Println(mysqlConfig)
	mCfg := mysql.Config{
		User:      mysqlConfig.Username,
		Passwd:    mysqlConfig.Password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%d", mysqlConfig.Host, mysqlConfig.Port),
		DBName:    mysqlConfig.Database,
		ParseTime: mysqlConfig.ParseTime,
		Loc:       time.UTC,
	}

	connector, err := mysql.NewConnector(&mCfg)
	if err != nil {
		log.Println("Failed to create MySQL connector.", err.Error())
		return nil, err
	}

	conn := sql.OpenDB(connector)
	if err := conn.Ping(); err != nil {
		log.Println("Failed to connect MySQL Server.", err.Error())
		return nil, err
	}
	log.Println("Connected to MySQL Server successfully.")
	// Set connection pool parameters
	conn.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	conn.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(mysqlConfig.ConnMaxLifetime))
	return conn, nil
}

// service start with http endpoint
func (s *server) Start() {
	//setup repository
	eventRepo := repository.NewEventRepository(s.mysqlDB)
	transactionManager := repository.NewTransactionManager(s.mysqlDB)

	//setup service
	eventService := service.NewEventService(transactionManager, eventRepo)

	//setup handler
	eventHandler := handler.NewEventHandler(eventService)

	//setup http server
	r := mux.NewRouter()
	//basic health api
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	//event related api
	r.HandleFunc("/events", eventHandler.InsertEvent).Methods(http.MethodPost)
	r.HandleFunc("/events", eventHandler.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{event_id}", eventHandler.DeleteEvent).Methods(http.MethodDelete)

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
