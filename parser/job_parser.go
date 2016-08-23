package parser

import (
	"errors"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/instructions"
)

type JobParser struct{}

func (j *JobParser) printError(err error) error {
	log.WithFields(log.Fields{
		"package": "parser",
	}).Error(err)

	return err
}

//CreateJobWithFilePaths will take a tasks and a infrastructure file as
//arguments to pair them and use the job dispatcher (also in this file)
func (j *JobParser) CreateJobWithFilePaths(tasksFilePath, infrastructureFilePath string) error {
	taskParser := new(TaskFileParser)
	taskReader, err := taskParser.Parse(tasksFilePath)
	if err != nil {
		return j.printError(err)
	}
	rawTaskList, err := taskParser.Build(taskReader)
	if err != nil {
		return j.printError(err)
	}
	taskList, err := j.ParseTaskList(rawTaskList)
	if err != nil {
		return j.printError(err)
	}

	infrastructureParser := new(InfrastructureFileParser)
	infrastructureReader, err := infrastructureParser.Parse(infrastructureFilePath)
	if err != nil {
		return j.printError(err)
	}
	infrastructure, err := infrastructureParser.Build(infrastructureReader)
	if err != nil {
		return j.printError(err)
	}

	jobs := j.BuildJobList(infrastructure, taskList)

	//Send jobs to dispatcher
	return dispatcher.Dispatch(jobs)
}

func (j *JobParser) BuildJobList(infrastructure *raccoon.Infrastructure, tasks *[]raccoon.Task) *[]raccoon.Job {
	jobs := make([]raccoon.Job, 0)

	for _, cluster := range infrastructure.Infrastructure {
		for _, task := range cluster.TasksToExecute {
			for _, commands := range *tasks {
				//Find the associated tasks for each cluster
				if strings.ToLower(commands.Title) == strings.ToLower(task) {
					jobs = append(jobs, raccoon.Job{
						Cluster: cluster,
						Task:    commands,
					})
				}
			}
		}
	}

	return &jobs
}

//ParseTaskList takes a taskList (group of commands) and check the
//commands of each instruction to assign the proper strategy
func (j *JobParser) ParseTaskList(rawTasks *[]raccoon.Task) (*[]raccoon.Task, error) {

	taskList := []raccoon.Task{}

	for _, rawTask := range *rawTasks {
		parsedInstructions := make([]raccoon.CommandsExecutor, 0)

		for _, i := range rawTask.Command {
			if i["name"] == "" {
				return nil, errors.New("No \"name\" found on instructions")
			}

			switch i["name"] {
			case "RUN":
				if i["instruction"] == "" {
					return nil, errors.New("At least one " +
						"'instruction' was missing on a RUN command in the " +
						"tasks file")
				}
				run := instructions.RUN{
					Command: raccoon.Command{
						Name:        "RUN",
						Description: i["description"],
					},
					Instruction: i["instruction"],
				}
				parsedInstructions = append(parsedInstructions, &run)
			case "ADD":
				if i["sourcePath"] == "" {
					return nil, errors.New("SourcePath not specified" +
						"on ADD instruction")
				}

				if i["destPath"] == "" {
					return nil, errors.New("destPath not specified" +
						"on ADD instruction")
				}
				add := instructions.ADD{
					Command: raccoon.Command{
						Name:        "RUN",
						Description: i["description"],
					},
					SourcePath:  i["sourcePath"],
					DestPath:    i["destPath"],
				}
				parsedInstructions = append(parsedInstructions, &add)
			}
		}

		taskList = append(taskList, raccoon.Task{
			Title:      rawTask.Title,
			Maintainer: rawTask.Maintainer,
			Commands:   parsedInstructions,
		})
	}

	return &taskList, nil
}
