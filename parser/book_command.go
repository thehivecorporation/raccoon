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
func ExecuteCommandsFile(commandsFilePath, infrastructureFilePath string) error {
	if commandsFilePath == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			constants.COMMANDS_NAME, constants.COMMANDS_NAME)
		log.Error(err)
		return err
	}

	if infrastructureFilePath == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			constants.HOSTS_FLAG_NAME, constants.COMMANDS_NAME)
		log.Error(err)
		return err
	}

	//Check for the zombiebook specified in -f flag
	commandsFile, err := readCommandFile(commandsFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	//Check for the mansion specified in the -m flag
	infrastructureFile, err := readInfrastructureFile(infrastructureFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	return generateJobs(infrastructureFile, commandsFile)
}

func generateJobs(infrastructure *Infrastructure, zbook job.CommandsFile) error {
	jobs := make([]job.Job, 0)
	for _, room := range infrastructure.Clusters {
		//Each room is a cluster
		for _, commands := range zbook {
			//Compare every assigned chapter to every cluster
			if strings.ToLower(commands.Title) == strings.ToLower(room.Commands) {
				jobs = append(jobs, job.Job{
					Cluster: connection.Cluster{
						Name:    room.Name,
						Commands: room.Commands,
						Nodes:   room.Hosts,
					},
					Commands: commands,
				})
			}
		}
	}

	//Send jobs to dispatcher
	return dispatcher.Dispatch(&jobs)
}
