package parser

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/job"
	"fmt"
)

//ExecuteZombieBook will take a zombiebook.json file and a mansion.file as
//arguments to pair them and use the job dispatcher (also in this file)
func ExecuteZombieBook(zbookFile, mansionFile string) error {
	if zbookFile == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			constants.INSTRUCTIONS_NAME, constants.INSTRUCTIONS_NAME)
		log.Error(err)
		return err
	}

	if mansionFile == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			constants.HOSTS_FLAG_NAME, constants.INSTRUCTIONS_NAME)
		log.Error(err)
		return err
	}

	//Check for the zombiebook specified in -f flag
	zbook, err := readZbookFile(zbookFile)
	if err != nil {
		log.Error(err)
		return err
	}

	//Check for the mansion specified in the -m flag
	mansion, err := readMansionFile(mansionFile)
	if err != nil {
		log.Error(err)
		return err
	}

	return generateJobs(mansion, zbook)
}

func generateJobs(mansion *mansion, zbook job.Zbook) error {
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
	return dispatcher.Dispatch(&jobs)
}
