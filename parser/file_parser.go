package parser

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

type FileParser struct{}

func (t *FileParser) Parse(filePath string) (io.Reader, error) {
	log.WithFields(log.Fields{
		"cluster": filePath,
		"package": "parser",
	}).Info("Reading " + filePath + " file")

	if filePath == "" {
		err := fmt.Errorf("You must provide a tasks and a infrastructure file. Check raccoon --help")
		log.Error(err)
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"tasks":   filePath,
			"package": "parser",
		}).Errorf("Could not read %s file\n", filePath)
		return nil, fmt.Errorf("Could not read %s file\n", filePath)
	}

	return file, nil
}
