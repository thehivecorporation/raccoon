package parser

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
)

//mansion is to parse json files and requests
type mansion struct {
	Name  string `json:"name"`
	Rooms []room `json:"rooms"`
}

//room is to parse json files and requests
type room struct {
	Name    string
	Chapter string
	Hosts   []connection.Node
}

//ReadMansionFile takes a filepath with a json containing a Mansion file and
//returns a Mansion file
func readMansionFile(f string) (*mansion, error) {
	log.WithFields(log.Fields{
		constants.HOST_NAME: f,
	}).Info(constants.ARROW_LENGTH + "Reading " + constants.HOSTS_FLAG_NAME + " file")

	var mansion_ mansion

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dat, &mansion_)
	if err != nil {
		return &mansion{}, err
	}

	return &mansion_, nil
}
