package raccoon

//Infrastructure represents a group of clusters. Each cluster represents a
//group of hosts.
type Infrastructure struct {
	//Name for your infrastructure. It has no effect on the app but maintains
	//infrastructure json files ordered
	Name           string    `json:"name"`

	//Infrastructure is, in turn, the array of clusters that you want to execute
	//some automation on
	Infrastructure []Cluster `json:"infrastructure"`
}
