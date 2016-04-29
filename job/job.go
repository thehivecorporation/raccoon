package job

import "github.com/thehivecorporation/raccoon/connection"

//Job is a entire task to perform on each Node of the "connection.Cluster" array
//and the recipe (list of tasks) that must be executed
type Job struct {
	Cluster *connection.Cluster
	Zbook *Zbook
}
