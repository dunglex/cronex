package main

import (
	"flag"
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	configFile := flag.String("config", "./cron.json", "Path to the cron configuration file")
	flag.Parse()

	jobs, err := ReadCronJobs(*configFile)
	if err != nil {
		fmt.Println("Could not read configuration file: ", err.Error())
	}

	cron := cron.New(cron.WithSeconds())
	for _, job := range jobs {
		if job.Enabled {
			if _, err := cron.AddJob(job.Cron, job); err != nil {
				panic(fmt.Sprintf("Could not add scheduled task: %s", err.Error()))
			}
			fmt.Println(job.ToString())
		}
	}
	cron.Run()
}
