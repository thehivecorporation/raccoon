package parser

import (
	"encoding/json"
	"io/ioutil"

	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
)

//mansion is to parse json files and requests
type Infrastructure struct {
	Name     string `json:"name"`
	Clusters []Cluster `json:"clusters"`
}

//room is to parse json files and requests
type Cluster struct {
	Name     string
	Commands string
	Hosts    []connection.Node
}

//readInfrastructureFile takes a file path with a json containing a Infrastructure file and
//returns a Infrastructure pointer
func readInfrastructureFile(f string) (*Infrastructure, error) {
	log.WithFields(log.Fields{
		constants.HOSTS_NAME: f,
		"package":            "parser",
	}).Info("Reading " + constants.HOSTS_FLAG_NAME + " file")

	var _inf Infrastructure

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dat, &_inf)
	if err != nil {
		return &Infrastructure{}, err
	}

	return checkErrors(&_inf)
}


//checkErrors is used to perform error checking on mansion json file
func checkErrors(m *Infrastructure) (*Infrastructure, error) {
	err := false
	if len(m.Clusters) == 0 {
		log.Error("No " + constants.HOSTS_NAME + " were found on " + constants.HOSTS_FLAG_NAME + " file")
		err = true
	}

	if m.Name == "" {
		log.Errorf("%s name can't be blank", constants.HOSTS_FLAG_NAME)
		err = true
	}

	for _, room := range m.Clusters {
		if len(room.Hosts) == 0 {
			log.Errorf("No hosts were found on %s '%s' for %s '%s'",
				constants.HOSTS_NAME, room.Name, constants.RELATIONSHIP_KEY,
				room.Commands)
			err = true
		}

		if room.Name == "" {
			log.Errorf("%s name can't be blank", constants.HOSTS_NAME)
			err = true
		}

		if room.Commands == "" {
			log.Errorf("%s name can't be blank on %s '%s'",
				constants.RELATIONSHIP_KEY, constants.HOSTS_NAME, room.Name)
			err = true
		}

		for _, host := range room.Hosts {
			if host.Username == "" {
				log.Errorf("Host username can't be blank on %s '%s'",
					constants.HOSTS_NAME, room.Name)
				err = true
			}

			if host.Password == "" {
				log.Errorf("Host password can't be blank on %s '%s'",
					constants.HOSTS_NAME, room.Name)
				err = true
			}

			if host.IP == "" {
				log.Errorf("Host IP can't be blank on %s '%s'",
					constants.HOSTS_NAME, room.Name)
				err = true
			}
		}
	}
	if err {
		return &Infrastructure{}, errors.New("Error found when parsing " + constants.HOSTS_FLAG_NAME + " file")
	}

	return m, nil
}
