package raccoon

//mansion is to parse json files and requests
type Infrastructure struct {
	Name           string    `json:"name"`
	Infrastructure []Cluster `json:"infrastructure"`
}
