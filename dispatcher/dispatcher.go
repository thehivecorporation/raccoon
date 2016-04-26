package dispatcher

import (
	"sync"

	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/job"
)

var wg sync.WaitGroup

//Dispatch receives a Job struct and is the one in charge of iterating for each
//Node within the job cluster of nodes and open a new routine for each of them
//with the recipe to execute.
func Dispatch(j job.Job) {
	for _, node := range j.Nodes {
		wg.Add(len(j.Recipe.Instructions))
		go executeRecipeOnNode(j.Recipe, node)
	}

	wg.Wait()
}

//executeRecipeOnNode will take every instruction of the recipe and execute it
//in order on each node. Each instruction waits until previous one is finished.
func executeRecipeOnNode(r job.Recipe, n connection.Node) {
	for _, instruction := range r.Instructions {
		instruction.Execute(n)
		wg.Done()
	}
}
