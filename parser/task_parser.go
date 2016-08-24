package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/thehivecorporation/raccoon"
)

//TaskFileParser is the type to build lists of tasks from an io.Reader
type TaskFileParser struct {
	FileParser
}

//Build takes an io.Reader with a JSON to parse it into a list of tasks
func (t *TaskFileParser) Build(r io.Reader) (*[]raccoon.Task, error) {
	var tasks []raccoon.Task
	err := json.NewDecoder(r).Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	return &tasks, nil
}
