package cronjobs

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/data"
)

func ExpireDataCronJob() {
	c, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("Error starting the Expire Data cron job. Error : ", err)
		return
	}

	job, err := c.NewJob(
		gocron.DurationJob(time.Duration(config.ExpireKeyCronTimer)*time.Second),
		gocron.NewTask(data.DeleteExpiredNodes),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Expire keys Cron Job Started.")
	log.Printf("Cron Info : %+v\n", job.ID())
	c.Start()
}

func BackupData() {
	c, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("Error starting the Backup Data cron job. Error : ", err)
		return
	}

	job, err := c.NewJob(
		gocron.DurationJob(time.Duration(config.BackupCronTimer)*time.Second),
		gocron.NewTask(data.BackupData),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Backup data Cron Job Started.")
	log.Printf("Cron Info : %+v\n", job.ID())
	c.Start()
}
