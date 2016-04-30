package parser

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/job"
)

//ExecuteZombieBook will take a zombiebook.json file and a mansion.file as
//arguments to pair them and use the job dispatcher (also in this file)
func ExecuteZombieBook(c *cli.Context) {
	if c.String(constants.INSTRUCTIONS_NAME) == "" {
		log.Fatalf("You must provide a %s file. Check raccoon %s --help",
			constants.INSTRUCTIONS_NAME, constants.INSTRUCTIONS_NAME)
	}

	if c.String(constants.HOSTS_FLAG_NAME) == "" {
		log.Fatalf("You must provide a %s file. Check raccoon %s --help",
			constants.HOSTS_FLAG_NAME, constants.INSTRUCTIONS_NAME)
	}

	//Check for the zombiebook specified in -f flag
	zbook, err := readZbookFile(c.String(constants.INSTRUCTIONS_NAME))
	if err != nil {
		log.Fatal(err)
	}

	//Check for the mansion specified in the -m flag
	mansion, err := readMansionFile(c.String(constants.HOSTS_FLAG_NAME))
	if err != nil {
		log.Fatal(err)
	}

	generateJobs(mansion, zbook)
}

func generateJobs(mansion *mansion, zbook job.Zbook) {
	jobs := make([]job.Job, 0)
	for _, room := range mansion.Rooms {
		//Each room is a cluster
		for _, chapter := range zbook {
			//Compare every assigned chapter to every cluster
			if strings.ToLower(chapter.Title) == strings.ToLower(room.Chapter) {
				jobs = append(jobs, job.Job{
					Cluster: connection.Cluster{
						Name:    room.Name,
						Chapter: room.Chapter,
						Nodes:   room.Hosts,
					},
					Chapter: chapter,
				})
			}
		}
	}

	//Send jobs to dispatcher
	dispatcher.Dispatch(&jobs)
}
