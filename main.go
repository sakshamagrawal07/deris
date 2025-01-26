package main

import (
	"flag"
	"log"
	"sync"

	"github.com/sakshamagrawal07/deris/data"
	"github.com/sakshamagrawal07/deris/server"
)

func setupFlags() {
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Host for the deris server")
	flag.IntVar(&server.Port, "port", 7379, "Port for the deris server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting deris...")
	// server.RunSyncTCPServer()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		data.ExpireDataCronJob()
	}()

	go func() {
		defer wg.Done()
		server.StartServer("localhost")
	}()

	wg.Wait()
}
