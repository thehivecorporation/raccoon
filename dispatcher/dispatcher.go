package dispatcher

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/job"
)

// A Sync group to wait all sessions to finish before exiting the app
var wg sync.WaitGroup

//Dispatch receives a Job struct and is the one in charge of iterating for each
//Node within the job cluster of nodes and open a new routine for each of them
//with the recipe to execute.
func Dispatch(js *[]job.Job) error {

	for _, job := range *js {
		log.WithFields(log.Fields{
			constants.HOSTS_NAME: job.Cluster.Name,
			constants.GROUP_NAME: job.Chapter.Title,
			constants.MAINTAINER: job.Chapter.Maintainer,
			"package":            "dispatcher",
		}).Info(constants.HOST_LAUNCH_MESSAGE)

		for _, node := range job.Cluster.Nodes {
			wg.Add(1)
			go executeRecipeOnNode(job, node)
		}
	}

	wg.Wait()
	return nil
}

//executeRecipeOnNode will take every instruction of the recipe and execute it
//in order on each node. Each instruction waits until previous one is finished.
func executeRecipeOnNode(j job.Job, n connection.Node) {
	n.GenerateUniqueColor()

	for _, instruction := range j.Chapter.Instructions {
		instruction.Execute(n)
	}

	err := n.CloseNode()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    n.IP,
			"package": "dispatcher",
		}).Warn("Error closing session: " + err.Error())
	}

	wg.Done()
}
