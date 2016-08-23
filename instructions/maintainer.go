package instructions

import "github.com/thehivecorporation/raccoon"

type MAINTAINER struct {
	Name        string
	Instruction string
	Description string
}

//Execute is the implementation of the Instruction interface for a MAINTAINER instruction TODO
func (m *MAINTAINER) Execute(h raccoon.Host) {
	logCommand(nil, h.IP, m.Description, m.GetCommandName())
}

func (m *MAINTAINER) GetCommandName() string {
	return "MAINTAINER"
}