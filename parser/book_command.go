package raccoon_cli

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/job"
	"github.com/thehivecorporation/raccoon/parser"
)

func executeZombieBook(c *cli.Context) {
	if c.String(constants.INSTRUCTIONS_NAME) == "" {
		log.Fatalf("You must provide a %s file. Check raccoon %s --help",
			constants.INSTRUCTIONS_NAME, constants.INSTRUCTIONS_NAME)
	}

	if c.String(constants.HOSTS_FLAG_NAME) == "" {
		log.Fatalf("You must provide a %s file. Check raccoon %s --help",
			constants.HOSTS_FLAG_NAME, constants.INSTRUCTIONS_NAME)
	}

	//Check for the zombiebook specified in -f flag
	zbook, err := parser.ReadZbookFile(c.String(constants.INSTRUCTIONS_NAME))
	if err != nil {
		log.Fatal(err)
	}

	//Check for the mansion specified in the -m flag
	mansion, err := parser.ReadMansionFile(c.String(constants.HOSTS_FLAG_NAME))
	if err != nil {
		log.Fatal(err)
	}

	//Generate Jobs, each job must be associated with a chapter title.
	jobs := make([]job.Job, 0)
	for _, room := range mansion.rooms {
		//Each room is a cluster
		for _, chapter := range zbook {
			//Compare every assigned chapter to every cluster
			if strings.ToLower(chapter.Title) == strings.ToLower(room.Chapter) {
				jobs = append(jobs, job.Job{
					Cluster: room,
					Chapter: chapter,
				})
			}
		}
	}

	//Send jobs to dispatcher
	dispatcher.Dispatch(&jobs)
}
