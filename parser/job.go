package parser

import (
	"errors"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/instructions"
)

//JobParser is used to parse JSON objects (tasks and infrastructure files) into their corresponding types
type Job struct {
	Generic
	Dispatcher raccoon.Dispatcher
}

func (j *Job) printError(err error) error {
	log.WithFields(log.Fields{
		"package": "parser",
	}).Error(err)

	return err
}

//CreateJobWithFilePaths will take a tasks and a infrastructure file as
//arguments to pair them and use the job dispatcher (also in this file)
func (j *Job) CreateJobWithFilePaths(tasksFilePath, infrastructureFilePath string) error {
	taskParser := new(Task)
	taskReader, err := taskParser.Parse(tasksFilePath)
	if err != nil {
		return j.printError(err)
	}
	var rawTaskList []raccoon.Task
	if err := taskParser.Build(taskReader, &rawTaskList); err != nil {
		return j.printError(err)
	}
	taskList, err := j.ParseTaskList(&rawTaskList)
	if err != nil {
		return j.printError(err)
	}

	infrastructureParser := new(InfrastructureFile)
	infrastructureReader, err := infrastructureParser.Parse(infrastructureFilePath)
	if err != nil {
		return j.printError(err)
	}
	var infrastructure raccoon.Infrastructure
	if err := infrastructureParser.Build(infrastructureReader, &infrastructure); err != nil {
		return j.printError(err)
	}
	infrastructureParser.Prepare(&infrastructure)

	jobs := j.BuildJobList(&infrastructure, taskList)

	//Send jobs to dispatcher
	j.Dispatcher.Dispatch(*jobs)

	return nil
}

//BuildJobList takes every cluster of an infrastructure file and searches its associated task in the task list to create the needed jobs
func (j *Job) BuildJobList(infrastructure *raccoon.Infrastructure, tasks *[]raccoon.Task) *[]raccoon.Job {
	jobs := make([]raccoon.Job, 0)

	for _, cluster := range infrastructure.Infrastructure {
		for _, taskToExecute := range cluster.TasksToExecute {
			for _, task := range *tasks {
				//Find the associated tasks for each cluster
				if strings.ToLower(task.Title) == strings.ToLower(taskToExecute) {
					jobs = append(jobs, raccoon.Job{
						Cluster: cluster,
						Task:    task,
					})
				}
			}
		}
	}

	return &jobs
}

//ParseTaskList takes a taskList (group of commands) and check the
//commands of each instruction to assign the proper strategy
func (j *Job) ParseTaskList(rawTasks *[]raccoon.Task) (*[]raccoon.Task, error) {

	taskList := []raccoon.Task{}

	for _, rawTask := range *rawTasks {
		parsedInstructions := make([]raccoon.CommandsExecutor, 0)

		for _, i := range rawTask.Command {
			if i["name"] == "" {
				return nil, errors.New("No 'name' found on instructions")
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
						Name:        "ADD",
						Description: i["description"],
					},
					SourcePath: i["sourcePath"],
					DestPath:   i["destPath"],
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
