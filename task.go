package raccoon

//Task is a list of commands with a name and a maintainer
type RawTask struct {
	Title      string              `json:"title"`
	Maintainer string              `json:"maintainer"`
	Command    []map[string]string `json:"command"`
}

//Task is a list of commands to execute on a host with a title and a maintainer
type Task struct {
	Title      string              `json:"title"`
	Maintainer string              `json:"maintainer"`
	Commands   []CommandsExecutor  `json:"commandList,omitempty"`
	Command    []map[string]string `json:"commands"`
}
