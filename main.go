package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/rahulshewale153/meeting-scheduler-api/configreader"
	"github.com/rahulshewale153/meeting-scheduler-api/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing config file \n usage: ./sigma-api-service <configPath>")
	}
	configFile := flag.String("config", "", "Service Configuration File")
	flag.Parse()

	config, err := configreader.ReadConfigFile(*configFile)
	if err != nil {
		log.Fatal("error while validating config file: %s", err.Error())
	}
	server := server.NewServer(config)

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server.Start()

	//waiting for interrupt or terminate signals
	<-stop
	log.Println("Shutting down...")
	server.Stop()

}
