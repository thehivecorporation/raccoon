package raccoon

//Task is a list of commands to execute on a host with a title and a maintainer
type Task struct {
	Title      string
	Maintainer string
	Commands   []CommandsExecutor
}
