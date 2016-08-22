package raccoon

//Job is that relate a Cluster with the Task that must be performed on it
type Job struct {
	Cluster Cluster
	Task    Task
}
