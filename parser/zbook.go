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
	//chapter is used to parse json files and requests
	chapter struct {
		Title           string              `json:"chapter_title"`
		Maintainer      string              `json:"maintainer"`
		RawInstructions []map[string]string `json:"instructions"`
	}

	//zombiebook is used to parse json files and requests
	zombiebook []chapter
)

//readZbookFile will take a filepath as parameter and return a Job
func readZbookFile(f string) (job.Zbook, error) {
	log.WithFields(log.Fields{
		constants.INSTRUCTIONS_NAME: f,
	}).Info(constants.ARROW_LENGTH + "Reading " + constants.INSTRUCTIONS_NAME +
		" file")

	var book zombiebook

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return job.Zbook{}, errors.New("Error reading " +
			constants.INSTRUCTIONS_FLAG_NAME + " file: " + err.Error())
	}

	err = json.Unmarshal(dat, &book)
	if err != nil {
		return job.Zbook{}, errors.New("Error parsing JSON: " + err.Error())
	}

	return generateZbookJob(book)
}

//generateZbookJob takes a zombiebook (group of instructions) and check the
//commands of each instruction to assign the proper strategy
func generateZbookJob(zombiebook zombiebook) (job.Zbook, error) {

	zbookJob := job.Zbook{}

	for _, z := range zombiebook {
		parsedInstructions := make([]instructions.InstructionExecutor, 0)

		for _, i := range z.RawInstructions {
			if i["name"] == "" {
				return job.Zbook{}, errors.New("No \"name\" found on instructions")
			}

			switch i["name"] {
			case "RUN":
				if i["instruction"] == "" {
					return job.Zbook{}, errors.New("At least one " +
						"'instruction' was missing on a RUN command in the " +
						constants.INSTRUCTIONS_NAME + " file")
				}
				run := instructions.RUN{
					Name:        "RUN",
					Description: i["description"],
					Instruction: i["instruction"],
				}
				parsedInstructions = append(parsedInstructions, &run)
			case "ADD":
				if i["sourcePath"] == "" {
					return job.Zbook{}, errors.New("SourcePath not specified" +
						"on ADD instruction")
				}

				if i["destPath"] == "" {
					return job.Zbook{}, errors.New("destPath not specified" +
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

		zbookJob = append(zbookJob, job.Chapter{
			Title:        z.Title,
			Maintainer:   z.Maintainer,
			Instructions: parsedInstructions,
		})
	}

	return zbookJob, nil
}
