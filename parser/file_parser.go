package parser

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

type FileParser struct{}

func (t *FileParser) Parse(filePath string) (io.Reader, error) {
	log.WithFields(log.Fields{
		raccoon.HOSTS_NAME: filePath,
		"package":  "parser",
	}).Info("Reading " + filePath + " file")

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
			"package":          "parser",
		}).Errorf("Could not read %s file\n", filePath)
		return nil, fmt.Errorf("Could not read %s file\n", filePath)
	}

	return file, nil
}
