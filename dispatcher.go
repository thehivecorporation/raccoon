package main

func Dispatch(j Job, q chan bool) {
	for _, node := range j.Nodes {
		go ExecuteRecipeOnNode(j.Recipe, node, q)
	}
}

func ExecuteRecipeOnNode(r Recipe, n Node, q chan bool) {
	for _, command := range r.Commands {
		ExecuteCommandOnNode(command, n)
	}

	q <- true
}
