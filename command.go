package raccoon

//Command is the type that contains common command members. Every new command in the commands
//package must have a Command value inside so that logging can work properly.
type Command struct {

	//Name if a must in each command. It could be ADD or ENV but every command
	//must have one so that the parser can recognize the intents of the user
	Name string

	//Description is an optional description of a command that will be print on
	//stdout for logging purposes. It's recommended to use a description on each
	//command
	Description string
}
