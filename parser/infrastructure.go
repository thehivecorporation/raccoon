package parser

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

//InfrastructureFileParser is the type that is in charge of parsing the files
//with the infrastructure information about a Raccoon job.
type InfrastructureFile struct {
	Generic
}

func (t *InfrastructureFile) Prepare(i *raccoon.Infrastructure) {
	t.TakeAuthAtClusterLevel(i)
	t.CheckErrors(i)
}

//TakeAuthAtClusterLevel will check if the user has written the authentication of a cluster in the
//cluster definition instead of on each individual host. In such case it will take the information
//and inject it on each host
func (t *InfrastructureFile) TakeAuthAtClusterLevel(i *raccoon.Infrastructure) {
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
}

//CheckErrors is used to perform error checking on Infrastructure json file
func (t *InfrastructureFile) CheckErrors(m *raccoon.Infrastructure) {
	if len(m.Infrastructure) == 0 {
		log.Error(parseErrorFactory(NO_CLUSTER))
	}

	if m.Name == "" {
		log.Errorf(parseErrorFactory(NO_CLUSTER_NAME).Error())
	}

	for _, cluster := range m.Infrastructure {
		if len(cluster.Hosts) == 0 {
			log.Errorf(parseErrorFactory(NO_HOSTS, cluster.Name, fmt.Sprintf("%#v", cluster.TasksToExecute)).Error())
		}

		if len(cluster.TasksToExecute) == 0 {
			log.Warnf(parseErrorFactory(NO_TASKS, cluster.Name).Error())
		}

		for _, host := range cluster.Hosts {
			if host.Username == "" {
				log.Errorf(parseErrorFactory(BLANK_USERNAME, host.IP).Error())
			}

			if host.Password == "" {
				log.Warnf(parseErrorFactory(BLANK_PASSWORD, host.IP).Error())
			}

			if host.IP == "" {
				log.Errorf(parseErrorFactory(BLANK_IP).Error())
			}
		}
	}
}
