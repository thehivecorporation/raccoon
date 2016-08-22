package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

type InfrastructureFileParser struct {
	FileParser
}

func (t *InfrastructureFileParser) Build(r io.Reader) (*raccoon.Infrastructure, error) {
	var infrastructure raccoon.Infrastructure
	err := json.NewDecoder(r).Decode(&infrastructure)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON: " + err.Error())
	}

	return t.checkErrors(&infrastructure)
}

//checkErrors is used to perform error checking on mansion json file
func (t *InfrastructureFileParser) checkErrors(m *raccoon.Infrastructure) (*raccoon.Infrastructure, error) {
	err := false
	if len(m.Infrastructure) == 0 {
		log.Error("No cluster were found on infrastructure file")
		err = true
	}

	if m.Name == "" {
		log.Errorf("infrastructure name can't be blank")
		err = true
	}

	for _, cluster := range m.Infrastructure {
		if len(cluster.Hosts) == 0 {
			log.Errorf("No hosts were found on cluster '%s' for commands '%s'",
				cluster.Name, cluster.Task)
			err = true
		}

		if cluster.Name == "" {
			log.Errorf("hosts name can't be blank")
			err = true
		}

		if cluster.Task == "" {
			log.Errorf("commands name can't be blank on cluster '%s'", cluster.Name)
			err = true
		}

		for _, host := range cluster.Hosts {
			if host.Username == "" {
				log.Errorf("Host username can't be blank on cluster '%s'", cluster.Name)
				err = true
			}

			if host.Password == "" {
				log.Errorf("Host password can't be blank on cluster '%s'", cluster.Name)
				err = true
			}

			if host.IP == "" {
				log.Errorf("Host IP can't be blank on cluster '%s'", cluster.Name)
				err = true
			}
		}
	}
	if err {
		return &raccoon.Infrastructure{}, errors.New("Error found when parsing " +
			"infrastructure file")
	}

	return m, nil
}
