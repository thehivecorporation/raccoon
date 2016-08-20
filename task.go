package raccoon

//Task is a list of commands with a name and a maintainer
type RawTask struct {
	Title      string              `json:"title"`
	Maintainer string              `json:"maintainer"`
	Command    []map[string]string `json:"command"`
}
