package main

import (
	"flag"
	"log"
	"sync"

	"github.com/sakshamagrawal07/deris/cronjobs"
	"github.com/sakshamagrawal07/deris/server"
)

func setupFlags() {
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Host for the deris server")
	flag.IntVar(&server.Port, "port", 7379, "Port for the deris server")
	flag.IntVar(&server.ExpireKeyCronTimer, "expire-key-cron-timer", 60, "Time in seconds for the expire keys cronjob.")
	flag.IntVar(&server.BackupCronTimer, "backup-cron-timer", 60, "Time in seconds for the backup data cronjob.")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting deris...")
	// server.RunSyncTCPServer()

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		cronjobs.ExpireDataCronJob()
	}()

	go func() {
		defer wg.Done()
		cronjobs.BackupData()
	}()

	go func() {
		defer wg.Done()
		server.StartServer("localhost")
	}()

	wg.Wait()
}
