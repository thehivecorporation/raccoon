package parser

import (
	"encoding/json"
	"io/ioutil"

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
		return job.Zbook{}, err
	}

	err = json.Unmarshal(dat, &book)
	if err != nil {
		return job.Zbook{}, err
	}

	return generateZbookJob(book), nil
}

//generateZbookJob takes a zombiebook (group of instructions) and check the
//commands of each instruction to assign the proper strategy
func generateZbookJob(zombiebook zombiebook) job.Zbook {

	zbookJob := job.Zbook{}

	for _, z := range zombiebook {
		parsedInstructions := make([]instructions.InstructionExecutor, 0)

		for _, i := range z.RawInstructions {
			switch i["name"] {
			case "RUN":
				run := instructions.RUN{
					Name:        "RUN",
					Description: i["description"],
					Instruction: i["instruction"],
				}
				parsedInstructions = append(parsedInstructions, &run)
			case "ADD":
				add := instructions.ADD{
					SourcePath:  i["sourcePath"],
					DestPath:    i["destPath"],
					Description: i["description"],
					Name:        "ADD",
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

	return zbookJob
}
