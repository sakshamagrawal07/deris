package main

import (
	"flag"
	"log"
	"sync"

	"github.com/sakshamagrawal07/deris/commands"
	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/cronjobs"
	"github.com/sakshamagrawal07/deris/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for the deris server")
	flag.IntVar(&config.Port, "port", 7379, "Port for the deris server")
	flag.IntVar(&config.ExpireKeyCronTimer, "expire-key-cron-timer", 60, "Time in seconds for the expire keys cronjob.")
	flag.IntVar(&config.BackupCronTimer, "backup-cron-timer", 60, "Time in seconds for the backup data cronjob.")
	flag.BoolVar(&config.ClearAOF, "clear-aof", true, "Clear the AOF file and start a fresh server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Starting deris...")
	// server.RunSyncTCPServer()

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		cronjobs.ExpireDataCronJob()
	}()

	go func() {
		defer wg.Done()
		commands.ExecuteCommandsInQueue()
	}()

	go func() {
		defer wg.Done()
		cronjobs.BackupData()
	}()

	go func() {
		defer wg.Done()
		server.InitData()
		if !config.ClearAOF {
			commands.RecoverFromAOF()
		} else {
			commands.ClearAof()
		}
		server.StartServer("localhost")
	}()

	wg.Wait()
}
