package raccoon

//Job is a entire instruction list to perform on each Node of the "connection.Cluster" array
//and the book (list of instructions) that must be executed
type Job struct {
	Cluster Cluster
	Task    Task
}
