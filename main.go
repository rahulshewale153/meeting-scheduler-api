package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/rahulshewale153/meeting-scheduler-api/server"
)

func main() {
	port := 8080
	server := server.NewServer(port)
	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server.Start()

	//waiting for interrupt or terminate signals
	<-stop
	log.Println("Shutting down...")
	server.Stop()

}
