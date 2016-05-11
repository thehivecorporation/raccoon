package instructions

import (
	log "github.com/Sirupsen/logrus"

	"github.com/thehivecorporation/raccoon/connection"
)

//RUN is a instruction that in the recipe file correspond to the CMD instruction.
//It will execute the "Command" on every machine. Ideally, every command must
//be bash
type RUN struct {
	Name        string
	Description string
	Instruction string
}

//Execute is the implementation of the Instruction interface for a CMD instruction
func (c *RUN) Execute(n connection.Node) {
	session, err := n.GetSession()
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "RUN",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error(err.Error())
	}

	log.WithFields(log.Fields{
		"Instruction": "RUN",
		"Node":        n.IP,
		"package":     "instructions",
	}).Info(c.Description)

	err = session.Run(c.Instruction)
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "RUN",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error("Error running command on session: " + err.Error())
	}

	session.Close()
}
