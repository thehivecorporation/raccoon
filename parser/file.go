package parser

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

//FileParser is the type that parses files from a file path string. It manages
//file opening and returns errors.
type File struct{}

//Parse takes a filePath and returns an io.Reader pointing to it so that the
//Task and Job parser can work on it
func (t *File) Parse(filePath string) (io.Reader, error) {
	log.WithFields(log.Fields{
		"cluster": filePath,
		"package": packageName,
	}).Info("Reading " + filePath + " file")

	if filePath == "" {
		err := fmt.Errorf("You must provide a tasks and a infrastructure file or a job file. Check raccoon --help")
		log.Error(err)
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"tasks":   filePath,
			"package": packageName,
		}).Errorf("Could not read %s file\n", filePath)

		return nil, fmt.Errorf("Could not read %s file\n", filePath)
	}

	return file, nil
}
