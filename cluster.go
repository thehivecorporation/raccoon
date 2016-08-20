package raccoon

//Cluster is an array of Nodes (remotes machines) that compose our machine cluster
type Cluster struct {
	Name     string `json:"name"`
	Hosts    []Host `json:"hosts"`
	Commands string `json:"commands"`
}
