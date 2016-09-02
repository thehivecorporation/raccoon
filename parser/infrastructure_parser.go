package parser

import (
	"encoding/json"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

//InfrastructureFileParser is the type that is in charge of parsing the files
//with the infrastructure information about a Raccoon job.
type InfrastructureFileParser struct {
	FileParser
}

//Build returns an Infrastructure from a io.Reader
func (t *InfrastructureFileParser) Build(r io.Reader) (*raccoon.Infrastructure, error) {
	var infrastructure raccoon.Infrastructure
	err := json.NewDecoder(r).Decode(&infrastructure)
	if err != nil {
		return nil, infrastructureErr(JSON_ERROR, err.Error())
	}

	t.TakeAuthAtClusterLevel(&infrastructure)

	return t.CheckErrors(&infrastructure)
}

//takeAuthAtClusterLevel will check if the user has written the authentication of a cluster in the
//cluster definition instead of on each individual host. In such case it will take the information
//and inject it on each host
func (t *InfrastructureFileParser) TakeAuthAtClusterLevel(i *raccoon.Infrastructure) error {
	for k, cluster := range i.Infrastructure {
		var username, authFile, password string
		var interactive bool

		if i.Infrastructure[k].Username != "" {
			username = i.Infrastructure[k].Username
		}

		if i.Infrastructure[k].Password != "" {
			password = i.Infrastructure[k].Password
		}

		if i.Infrastructure[k].IdentityFile != "" {
			authFile = i.Infrastructure[k].IdentityFile
		}

		if i.Infrastructure[k].InteractiveAuth {
			interactive = i.Infrastructure[k].InteractiveAuth
		}

		for j := range cluster.Hosts {
			if i.Infrastructure[k].Hosts[j].Username == "" {
				i.Infrastructure[k].Hosts[j].Username = username
			}

			if i.Infrastructure[k].Hosts[j].Password == "" {
				i.Infrastructure[k].Hosts[j].Password = password
			}

			if i.Infrastructure[k].Hosts[j].IdentityFile == "" {
				i.Infrastructure[k].Hosts[j].IdentityFile = authFile
			}

			if !i.Infrastructure[k].Hosts[j].InteractiveAuth {
				i.Infrastructure[k].Hosts[j].InteractiveAuth = interactive
			}
		}
	}

	return nil
}

//checkErrors is used to perform error checking on Infrastructure json file
func (t *InfrastructureFileParser) CheckErrors(m *raccoon.Infrastructure) (*raccoon.Infrastructure, error) {
	err := false
	if len(m.Infrastructure) == 0 {
		log.Error(infrastructureErr(NO_CLUSTER))
		err = true
	}

	if m.Name == "" {
		log.Errorf(infrastructureErr(NO_CLUSTER_NAME).Error())
		err = true
	}

	for _, cluster := range m.Infrastructure {
		if len(cluster.Hosts) == 0 {
			log.Errorf(infrastructureErr(NO_HOSTS, cluster.Name, fmt.Sprintf("%#v",cluster.TasksToExecute)).Error())
			err = true
		}

		if len(cluster.TasksToExecute) == 0 {
			log.Errorf(infrastructureErr(NO_TASKS, cluster.Name).Error())
			err = true
		}

		for _, host := range cluster.Hosts {
			if host.Username == "" {
				log.Errorf(infrastructureErr(BLANK_USERNAME, host.IP).Error())
				err = true
			}

			if host.Password == "" {
				log.Warnf(infrastructureErr(BLANK_PASSWORD, host.IP).Error())
			}

			if host.IP == "" {
				log.Errorf(infrastructureErr(BLANK_IP).Error())
				err = true
			}
		}
	}
	if err {
		return &raccoon.Infrastructure{}, infrastructureErr(PARSING_ERROR)
	}

	return m, nil
}

//Errors
const (
	PARSING_ERROR  string = "Error found when parsing infrastructure file"
	BLANK_IP       string = "Host IP can't be blank"
	BLANK_PASSWORD string = "Host password is blank on host '%s'. If no password is specified you " +
		"must use an identity file or an interactive authentication method"
	BLANK_USERNAME string = "Host username is blank on host '%s'"
	NO_TASKS       string = "You haven't specified any task. Specify at least one as an string " +
		"array on cluster '%s'"
	NO_HOSTS        string = "No hosts were found on cluster '%s' for commands '%s'"
	NO_CLUSTER_NAME string = "infrastructure name can't be blank"
	NO_CLUSTER      string = "No cluster was found on infrastructure file"
	JSON_ERROR      string = "Error parsing JSON: %s\n"
)

type InfrastructureError struct {
	errorCode string
	msg       string
	extra     []string
}

func (i *InfrastructureError) Error() string {
	return fmt.Sprintf(i.msg, i.extra)
}

func infrastructureErr(errorCode string, extra ...string) error {
	err := InfrastructureError{
		errorCode: errorCode,
		msg:       errorCode,
		extra:     extra,
	}

	return &err
}
