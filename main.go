package main

import (
	"flag"
	"log"

	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for the deris server")
	flag.IntVar(&config.Port, "port", 7379, "Port for the deris server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting deris...")
	// server.RunSyncTCPServer()
	server.StartServer("localhost")
}
