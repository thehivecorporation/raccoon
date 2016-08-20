package raccoon

import (
	"encoding/json"
	"io/ioutil"

	"errors"

	log "github.com/Sirupsen/logrus"
)

//mansion is to parse json files and requests
type Infrastructure struct {
	Name     string    `json:"name"`
	Clusters []Cluster `json:"clusters"`
}

//readInfrastructureFile takes a file path with a json containing a Infrastructure file and
//returns a Infrastructure pointer
func ReadInfrastructureFile(f string) (*Infrastructure, error) {
	log.WithFields(log.Fields{
		HOSTS_NAME: f,
		"package":  "parser",
	}).Info("Reading " + INFRASTRUCTURE + " file")

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
		log.Error("No " + HOSTS_NAME + " were found on " + INFRASTRUCTURE + " file")
		err = true
	}

	if m.Name == "" {
		log.Errorf("%s name can't be blank", INFRASTRUCTURE)
		err = true
	}

	for _, cluster := range m.Clusters {
		if len(cluster.Hosts) == 0 {
			log.Errorf("No hosts were found on %s '%s' for %s '%s'",
				HOSTS_NAME, cluster.Name, RELATIONSHIP_KEY,
				cluster.Commands)
			err = true
		}

		if cluster.Name == "" {
			log.Errorf("%s name can't be blank", HOSTS_NAME)
			err = true
		}

		if cluster.Commands == "" {
			log.Errorf("%s name can't be blank on %s '%s'",
				RELATIONSHIP_KEY, HOSTS_NAME, cluster.Name)
			err = true
		}

		for _, host := range cluster.Hosts {
			if host.Username == "" {
				log.Errorf("Host username can't be blank on %s '%s'",
					HOSTS_NAME, cluster.Name)
				err = true
			}

			if host.Password == "" {
				log.Errorf("Host password can't be blank on %s '%s'",
					HOSTS_NAME, cluster.Name)
				err = true
			}

			if host.IP == "" {
				log.Errorf("Host IP can't be blank on %s '%s'",
					HOSTS_NAME, cluster.Name)
				err = true
			}
		}
	}
	if err {
		return &Infrastructure{}, errors.New("Error found when parsing " + INFRASTRUCTURE + " file")
	}

	return m, nil
}
