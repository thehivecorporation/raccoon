package instructions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/tmc/scp"
	"github.com/thehivecorporation/raccoon/connection"
)

//ADD is a instruction to copy files with scp
type ADD struct {
	//Source path of the file
	SourcePath  string
	//Destination path of the file
	DestPath    string
	//Description of the instruction
	Description string
	//The name that identifies this struct ("ADD" in this case)
	Name        string
}

//Execute is the implementation of the Instruction interface for a ADD instruction
func (c *ADD) Execute(n connection.Node) {
	session,err := n.GetSession()
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "RUN",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error(err.Error())
		session.Close()

		return
	}

	log.WithFields(log.Fields{
		"Instruction": "ADD",
		"Node":        n.IP,
		"SourcePath":  c.SourcePath,
		"DestPath":    c.DestPath,
		"package":     "instructions",
	}).Info(c.Description)

	scp.CopyPath(c.SourcePath,c.DestPath, session)

	session.Close()
}
