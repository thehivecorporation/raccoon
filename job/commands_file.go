package job

import "github.com/thehivecorporation/raccoon/instructions"

//CommandsFile is a list of Commands that must be executed, in order, on a remote
//machine
type CommandsFile []CommandsList

type CommandsList struct {
	Title      string
	Maintainer string
	Commands   []instructions.InstructionExecutor
}