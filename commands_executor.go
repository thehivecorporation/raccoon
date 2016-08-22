package raccoon

//CommandsExecutor is an interface that every command must implement. A command
//is any strategy of the Dockerfile syntax that Raccoon offers.
type CommandsExecutor interface {

	//Execute is the method that each strategy will execute on the provided
	//Host n
	Execute(n Host)
}
