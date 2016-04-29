package raccoon_cli

import (
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

	//Generate jobs
	jobs := make([]job.Job, len(mansion.Rooms))
	for i, c := range mansion.Rooms {
		jobs[i] = job.Job{
			GroupName: mansion.Name,
			Cluster: &c,
			Zbook:   &zbook,
		}
	}

	//Send jobs to dispatcher
	dispatcher.Dispatch(&jobs)
}
