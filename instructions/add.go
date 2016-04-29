package instructions

import (
	"fmt"
	log "github.com/Sirupsen/logrus"

	"github.com/thehivecorporation/raccoon/connection"
)

type ADD struct {
	SourcePath  string
	DestPath    string
	Instruction string
}

//Execute is the implementation of the Instruction interface for a ADD instruction TODO
func (c *ADD) Execute(n connection.Node) {
	log.WithFields(log.Fields{
		"Instruction": "ADD",
		"Node":        n.IP,
	}).Info(fmt.Sprintf("------------------------------> ADD: %s on %s\n", c.SourcePath, c.DestPath))
}
