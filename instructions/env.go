package instructions

import (
	"fmt"
	"strings"

	"github.com/thehivecorporation/raccoon"
)

type ENV struct {
	Name        string
	Description string
	Environment string
}

func (e *ENV) GetCommandName() string {
	return "ENV"
}

func (e *ENV) Execute(h raccoon.Host) {
	session, err := h.GetSession()
	if err != nil {
		logError(err, e, &h)

		session.Close()
		return
	}

	logCommand(nil, h.IP, e.Description, e.GetCommandName())

	env := strings.Split(e.Environment, "=")
	if len(env) == 2 {
		session.Setenv(env[0], env[1])
	} else {
		logError(fmt.Errorf("More than one '=' found on ENV instruction"), e, &h)
	}
}
