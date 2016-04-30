package job

import "github.com/thehivecorporation/raccoon/connection"

//Job is a entire instruction list to perform on each Node of the "connection.Cluster" array
//and the book (list of instructions) that must be executed
type Job struct {
	Cluster connection.Cluster
	Chapter Chapter
}
