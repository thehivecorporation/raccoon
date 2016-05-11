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
		constants.HOSTS_NAME: f,
		"package":"parser",
	}).Info("Reading " + constants.HOSTS_FLAG_NAME + " file")

	var mansion_ mansion

	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dat, &mansion_)
	if err != nil {
		return &mansion{}, err
	}

	return checkErrors(&mansion_)
}

func checkErrors(m *mansion) (*mansion, error) {
	err := false
	if len(m.Rooms) == 0 {
		log.Error("No rooms were found on " + constants.HOSTS_FLAG_NAME + " file")
		err = true
	}

	if m.Name == "" {
		log.Errorf("%s name can't be blank", constants.HOSTS_FLAG_NAME)
		err = true
	}

	for _, room := range m.Rooms {
		if len(room.Hosts) == 0 {
			log.Errorf("No hosts were found on %s '%s' for %s '%s'",
				constants.HOSTS_NAME, room.Name, constants.RELATIONSHIP_KEY,
				room.Chapter)
			err = true
		}

		if room.Name == "" {
			log.Errorf("%s name can't be blank", constants.HOSTS_NAME)
			err = true
		}

		if room.Chapter == "" {
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
		return &mansion{}, errors.New("Error found when parsing " + constants.HOSTS_FLAG_NAME + " file")
	}

	return m, nil
}
