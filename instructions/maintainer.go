package instructions

import "github.com/thehivecorporation/raccoon"

type MAINTAINER struct {
	Command raccoon.Command
}

//Execute is the implementation of the Instruction interface for a MAINTAINER instruction TODO
func (m *MAINTAINER) Execute(h raccoon.Host) {
	logCommand(nil, h, m)
}

func (m *MAINTAINER) GetCommand() *raccoon.Command {
	return &m.Command
}