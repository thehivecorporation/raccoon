package instructions

import (
	"fmt"

	"github.com/thehivecorporation/raccoon/connection"
)

type MAINTAINER struct {
	Name        string
	Instruction        string
}

//Execute is the implementation of the Instruction interface for a MAINTAINER instruction TODO
func (c *MAINTAINER) Execute(n connection.Node) {
	fmt.Printf("Maintainer: %s\n", c.Name)
}
