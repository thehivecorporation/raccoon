package instructions

import (
	log "github.com/Sirupsen/logrus"

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

	e.LogCommand(&h)

	env := strings.Split(e.Environment, "=")
	if len(env) == 2 {
		session.Setenv(env[0], env[1])
	} else {
		logError(fmt.Errorf("More than one '=' found on ENV instruction"), e, &h)
	}
}

func (e *ENV) LogCommand(h *raccoon.Host) {
	log.WithFields(log.Fields{
		"Instruction": e.GetCommandName(),
		"Node":        h.IP,
		"package":     packageName,
	}).Info(e.Description)
}
