package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"strings"
	"github.com/thehivecorporation/raccoon/dispatcher"
)

type RaccoonFileParser interface {
	Parse(string) (io.Reader, error)
}

type FileJobParser struct {
	FilePath string
}

type FileParser struct{}

func (t *FileParser) Parse(filePath string) (io.Reader, error) {
	if filePath == "" {
		err := fmt.Errorf("You must provide a %s and a %s file. Check raccoon --help",
			raccoon.TASKS_NAME, raccoon.INFRASTRUCTURE)
		log.Error(err)
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			raccoon.TASKS_NAME: filePath,
			"package":             "parser",
		}).Errorf("Could not read %s file\n", filePath)
		return nil, fmt.Errorf("Could not read %s file\n", filePath)
	}

	return file, nil
}

type RaccoonJobBuilder interface {
	Build(io.Reader) (interface{}, error)
}

type TaskFileParser struct {
	FilePath string
	FileParser
}

func (t *TaskFileParser) Build(r io.Reader) (*[]raccoon.RawTask, error) {
	var tasks []raccoon.RawTask
	err := json.NewDecoder(r).Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	return &tasks, nil
}

type InfrastructureFileParser struct {
	FilePath string
	FileParser
}

func (t *InfrastructureFileParser) Build(r io.Reader) (*raccoon.Infrastructure, error) {
	var infrastructure raccoon.Infrastructure
	err := json.NewDecoder(r).Decode(&infrastructure)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	return &infrastructure, nil
}

//CreateJobWithFiles will take a tasks and a infrastructure file as
//arguments to pair them and use the job dispatcher (also in this file)
func CreateJobWithFiles(tasksFilePath, infrastructureFilePath string) error {
	if tasksFilePath == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			raccoon.TASKS_NAME, raccoon.TASKS_NAME)
		log.Error(err)
		return err
	}

	if infrastructureFilePath == "" {
		err := fmt.Errorf("You must provide a %s file. Check raccoon %s --help",
			raccoon.INFRASTRUCTURE, raccoon.TASKS_NAME)
		log.Error(err)
		return err
	}

	//Check for the zombiebook specified in -f flag
	commandsFile, err := readTasksFile(tasksFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	//Check for the mansion specified in the -m flag
	infrastructureFile, err := raccoon.ReadInfrastructureFile(infrastructureFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	return generateJobs(infrastructureFile, commandsFile)
}

func generateJobs(infrastructure *raccoon.Infrastructure, tasks []raccoon.Task) error {
	jobs := make([]raccoon.Job, 0)
	for _, room := range infrastructure.Clusters {
		//Each room is a cluster
		for _, commands := range tasks {
			//Compare every assigned chapter to every cluster
			if strings.ToLower(commands.Title) == strings.ToLower(room.Commands) {
				jobs = append(jobs, raccoon.Job{
					Cluster: raccoon.Cluster{
						Name:     room.Name,
						Commands: room.Commands,
						Hosts:    room.Hosts,
					},
					Task: commands,
				})
			}
		}
	}

	//Send jobs to dispatcher
	return dispatcher.Dispatch(&jobs)
}
