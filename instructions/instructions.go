package instructions

import "github.com/thehivecorporation/raccoon/connection"

//Instruction is an interface that every instruction must implement according
//to a Strategy design pattern
type InstructionExecutor interface {
	Execute(n connection.Node)
}
