package raccoon

//Command is the type that contains common command members. Every new command in the commands
//package must have a Command value inside so that logging can work properly.
type Command struct {
	Name        string
	Description string
}