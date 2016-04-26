package job

import "github.com/thehivecorporation/raccoon/instructions"

//Recipe is a list of instructions that must be executed, in order, on a remote
//machine
type Recipe struct {
	Instructions []instructions.Instruction
}
