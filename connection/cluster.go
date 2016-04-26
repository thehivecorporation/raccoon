package connection

//Cluster is an array of Nodes (remotes machines) that compose our machine cluster
type Cluster struct {
	Nodes []Node
}
