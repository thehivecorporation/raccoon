package instructions

import "github.com/thehivecorporation/raccoon"

//RUN is a instruction that in the recipe file correspond to the CMD instruction.
//It will execute the "Command" on every machine. Ideally, every command must
//be bash
type RUN struct {
	Command raccoon.Command

	//Bash instruction to execute
	Instruction string
}

func (r *RUN) GetCommand() *raccoon.Command {
	return &r.Command
}

//Execute is the implementation of the Instruction interface for a RUN instruction
func (r *RUN) Execute(h raccoon.Host) {
	session, err := h.GetSession()
	if err != nil {
		logError(err, r, &h)
		return
	}
	defer session.Close()

	logCommand(nil, h, r)

	if err = session.Run(r.Instruction); err != nil {
		logError(err, r, &h)
	}
}
