package dispatcher

import (
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/job"
	"sync"
)

var Wg sync.WaitGroup

func Dispatch(j job.Job) {
	for _, node := range j.Nodes {
		Wg.Add(len(j.Recipe.Instructions))
		go executeRecipeOnNode(j.Recipe, node)
	}

	Wg.Wait()
}

func executeRecipeOnNode(r job.Recipe, n connection.Node) {
	for _, instruction := range r.Instructions {
		instruction.Execute(n)
		Wg.Done()
	}
}
