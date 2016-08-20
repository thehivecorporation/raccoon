package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/thehivecorporation/raccoon"
)

type TaskFileParser struct {
	FileParser
}

func (t *TaskFileParser) Build(r io.Reader) (*[]raccoon.Task, error) {
	var tasks []raccoon.Task
	err := json.NewDecoder(r).Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	return &tasks, nil
}
