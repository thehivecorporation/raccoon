package instructions

import (
	log "github.com/Sirupsen/logrus"

	"github.com/thehivecorporation/raccoon"
)

//RUN is a instruction that in the recipe file correspond to the CMD instruction.
//It will execute the "Command" on every machine. Ideally, every command must
//be bash
type RUN struct {
	//The name that identifies this struct ("RUN" in this case)
	Name string

	//Description of the instruction that must be set by the user
	Description string

	//Bash instruction to execute
	Instruction string
}

//Execute is the implementation of the Instruction interface for a RUN instruction
func (c *RUN) Execute(n raccoon.Node) {
	session, err := n.GetSession()

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
		"Instruction": c.Name,
		"Node":        n.IP,
		"package":     "instructions",
	}).Info(c.Description)

	session.Run(c.Instruction)

	session.Close()
}
