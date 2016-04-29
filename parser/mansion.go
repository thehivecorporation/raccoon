package parser

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
)

type mansion struct {
	Name  string               `json:"mansion_name"`
	Rooms []connection.Cluster `json:"rooms"`
}

func ReadMansionFile(f string) (*[]connection.Cluster, error) {

	log.WithFields(log.Fields{
		constants.HOST_NAME: f,
	}).Info("------------------------------> Reading " + constants.HOSTS_FLAG_NAME + " file")

	var mansion mansion

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dat, &mansion)
	if err != nil {
		return nil, err
	}

	return &mansion.Rooms, nil
}
