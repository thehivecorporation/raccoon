package parser

import (
	"encoding/json"

	"errors"

	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/instructions"
)

//readTasksFile will take a filepath as parameter and return a Job
func readTasksFile(f string) ([]raccoon.Task, error) {
	file, err := os.Open(f)
	if err != nil {
		log.WithFields(log.Fields{
			raccoon.TASKS_NAME: f,
			"package":          "parser",
		}).Errorf("Could not read %s file\n", f)
		return []raccoon.Task{}, fmt.Errorf("Could not read %s file\n", f)
	}

	var tasks []raccoon.RawTask
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return []raccoon.Task{}, errors.New("Error parsing JSON: " + err.Error())
	}

	return generateCommandsJob(tasks)
}

//generateCommandsJob takes a taskList (group of commands) and check the
//commands of each instruction to assign the proper strategy
func generateCommandsJob(rawTasks []raccoon.RawTask) ([]raccoon.Task, error) {

	taskList := []raccoon.Task{}

	for _, rawTask := range rawTasks {
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
						raccoon.TASKS_NAME + " file")
				}
				run := instructions.RUN{
					Name:        "RUN",
					Description: i["description"],
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
					Name:        "ADD",
					SourcePath:  i["sourcePath"],
					DestPath:    i["destPath"],
					Description: i["description"],
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

	return taskList, nil
}
