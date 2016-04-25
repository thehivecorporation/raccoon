package instructions

import "github.com/thehivecorporation/raccoon/connection"

type Instruction interface {
	Execute(n connection.Node)
}
