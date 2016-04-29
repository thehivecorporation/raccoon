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

func ReadMansionFile(f string) (*mansion, error) {

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
		//Maybe is a single group file
		var singleGroup []connection.Node
		err = json.Unmarshal(dat, &singleGroup)
		if err != nil{
			return nil, err
		}

		cluster := make([]connection.Cluster,1)
		cluster[0] = connection.Cluster{
			Name:"main room",
			Nodes: singleGroup,
		}
		mansion_ = mansion{
			Name: "Apartment",
			Rooms: cluster,
		}

		return &mansion_, nil
	}

	return &mansion_, nil
}
