package instructions

import (
	log "github.com/Sirupsen/logrus"

	"github.com/thehivecorporation/raccoon/connection"
)

type ADD struct {
	SourcePath  string
	DestPath    string
	Description string
	Name        string
}

//Execute is the implementation of the Instruction interface for a ADD instruction TODO
func (c *ADD) Execute(n connection.Node) {
	log.WithFields(log.Fields{
		"Instruction": "ADD",
		"Node":        n.IP,
		"SourcePath":  c.SourcePath,
		"DestPath":    c.DestPath,
	}).Info("------------------------------> " + c.Description)
}
