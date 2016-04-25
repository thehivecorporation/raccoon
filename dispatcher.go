package main

func Dispatch(j Job) {
	for _, node := range j.Nodes {
		go ExecuteRecipeOnNode(j.Recipe, node)
	}
}

func ExecuteRecipeOnNode(r Recipe, n Node) {
	for _, command := range r.Commands {
		wg.Add(1)
		ExecuteCommandOnNode(command, n)
	}

}
