package instructions

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

type MAINTAINER struct {
	Name        string
	Instruction string
}

//Execute is the implementation of the Instruction interface for a MAINTAINER instruction TODO
func (m *MAINTAINER) Execute(n raccoon.Host) {
	fmt.Printf("Maintainer: %s\n", m.GetCommandName())
}

func (m *MAINTAINER) GetCommandName() string {
	return "MAINTAINER"
}

func (m *MAINTAINER) LogCommand(h *raccoon.Host) {
	log.WithFields(log.Fields{
		"Instruction": m.GetCommandName(),
		"Node":        h.IP,
		"package":     packageName,
	}).Info(m.Instruction)
}
