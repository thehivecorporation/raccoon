package raccoon

//Task is a list of commands to execute on a host. A task must have a Title
//that must match with the title that was written in the task key of the
//infrastructure file
type Task struct {

	//Title of the task that will be referenced from clusters
	Title      string              `json:"title"`

	//Maintainer is an optional member
	Maintainer string              `json:"maintainer,omitempty"`

	//Commands are the array of commands that will be executed on some host.
	//This member must never come from JSON that's why it has that key name
	Commands   []CommandsExecutor  `json:"noJson,omitempty"`

	//Command is the syntax definition of a command that Raccoon will interpret
	//into a CommandExecutor value
	Command    []map[string]string `json:"commands"`
}
