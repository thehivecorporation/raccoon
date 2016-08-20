package raccoon

//Instruction is an interface that every instruction must implement according
//to a Strategy design pattern
type CommandsExecutor interface {
	Execute(n Host)
}
