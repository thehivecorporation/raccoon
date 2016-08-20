package parser

import (
	"encoding/json"
	"github.com/thehivecorporation/raccoon"
)



//ParseRequest is called by the REST server to parse the JSON of a request
func ParseRequest(byt []byte) error {
	req := raccoon.Request{}
	err := json.Unmarshal(byt, &req)
	if err != nil {
		return err
	}

	generatedJobs, err := generateCommandsJob(req.CommandsFile)
	if err != nil {
		return err
	}

	generateJobs(&req.Infrastructure, generatedJobs)

	return nil
}
