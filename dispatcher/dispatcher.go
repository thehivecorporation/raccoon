package dispatcher

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

// A Sync group to wait all sessions to finish before exiting the app
var wg sync.WaitGroup

//Dispatch receives a Job struct and is the one in charge of iterating for each
//Node within the job cluster of nodes and open a new routine for each of them
//with the recipe to execute.
func Dispatch(js *[]raccoon.Job) error {

	for _, job := range *js {
		log.WithFields(log.Fields{
			"cluster": job.Cluster.Name,
			"infrastructure": job.Task.Title,
			"maintainer": job.Task.Maintainer,
			"package":          "dispatcher",
		}).Info("Launching Raccoon...")

		for _, node := range job.Cluster.Hosts {
			wg.Add(1)
			go executeRecipeOnNode(job, node)
		}
	}

	wg.Wait()
	return nil
}

//executeRecipeOnNode will take every instruction of the recipe and execute it
//in order on each node. Instructions are executed sequentially
func executeRecipeOnNode(j raccoon.Job, n raccoon.Host) {
	err := n.InitializeNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    n.IP,
			"package": "dispatcher",
		}).Warn("Error initializing node: " + err.Error())
	}

	for _, instruction := range j.Task.Commands {
		instruction.Execute(n)
	}

	err = n.CloseNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    n.IP,
			"package": "dispatcher",
		}).Warn("Error closing session: " + err.Error())
	}

	wg.Done()
}
