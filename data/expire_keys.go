package data

import (
	"log"

	"github.com/robfig/cron/v3"
)

func ExpireDataCronJob() {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", deleteExpiredNodes)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Cron Job Started.")
	log.Printf("Cron Info : %+v\n", c.Entries())
	c.Start()
}
