package raccoon

//Cluster is an array of Nodes (remotes machines) that compose our machine cluster
//
//It must have a name to identify it from other clusters.
//
//It should also have a linked task. A task is a group of commands that will be
//executed on this cluster
type Cluster struct {
	//Name that identifies this cluster respect to others
	Name           string `json:"name"`

	//Hosts are the array of hosts on this cluster
	Hosts          []Host `json:"hosts"`

	//TasksToExecute is the name of the commands group that will be executed in
	// this cluster. This name must match the name written in the 'Title' member
	// of the Tasks file. In future version it will work with an array os tasks too.
	TasksToExecute []string `json:"tasks"`
}
