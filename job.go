package raccoon

//Job is that relate a Cluster (group of hosts) with the Task that must be performed on it
type Job struct {

	//Cluster is a group of hosts
	Cluster Cluster

	//Task is a single task that could be composed of one or more commands
	Task Task
}
