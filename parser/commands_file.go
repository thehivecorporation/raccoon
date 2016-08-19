package parser

import (
	"encoding/json"
	"io/ioutil"

	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/instructions"
	"github.com/thehivecorporation/raccoon/job"
)

type (
	//Commands is used to parse json files and requests
	Commands struct {
		Title      string              `json:"title"`
		Maintainer string              `json:"maintainer"`
		Command    []map[string]string `json:"command"`
	}

	//commands is used to parse json files and requests
	commands []Commands
)

//readCommandFile will take a filepath as parameter and return a Job
func readCommandFile(f string) (job.CommandsFile, error) {
	log.WithFields(log.Fields{
		constants.COMMANDS_NAME: f,
		"package":"parser",
	}).Info("Reading " + constants.COMMANDS_NAME +
		" file")

	var commands commands

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return job.CommandsFile{}, errors.New("Error reading " +
			constants.COMMANDS_NAME + " file: " + err.Error())
	}

	err = json.Unmarshal(dat, &commands)
	if err != nil {
		return job.CommandsFile{}, errors.New("Error parsing JSON: " + err.Error())
	}

	return generateCommandsJob(commands)
}

//generateCommandsJob takes a zombiebook (group of instructions) and check the
//commands of each instruction to assign the proper strategy
func generateCommandsJob(commandsList commands) (job.CommandsFile, error) {

	commandsFile := job.CommandsFile{}

	for _, command := range commandsList {
		parsedInstructions := make([]instructions.InstructionExecutor, 0)

		for _, i := range command.Command {
			if i["name"] == "" {
				return job.CommandsFile{}, errors.New("No \"name\" found on instructions")
			}

			switch i["name"] {
			case "RUN":
				if i["instruction"] == "" {
					return job.CommandsFile{}, errors.New("At least one " +
						"'instruction' was missing on a RUN command in the " +
						constants.COMMANDS_NAME + " file")
				}
				run := instructions.RUN{
					Name:        "RUN",
					Description: i["description"],
					Instruction: i["instruction"],
				}
				parsedInstructions = append(parsedInstructions, &run)
			case "ADD":
				if i["sourcePath"] == "" {
					return job.CommandsFile{}, errors.New("SourcePath not specified" +
						"on ADD instruction")
				}

				if i["destPath"] == "" {
					return job.CommandsFile{}, errors.New("destPath not specified" +
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

		commandsFile = append(commandsFile, job.CommandsList{
			Title:        command.Title,
			Maintainer:   command.Maintainer,
			Commands: parsedInstructions,
		})
	}

	return commandsFile, nil
}
