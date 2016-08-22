package instructions

import (
	"github.com/thehivecorporation/raccoon"
	log "github.com/Sirupsen/logrus"
)

func logError(err error, c raccoon.CommandsExecutor, h *raccoon.Host){
	log.WithFields(log.Fields{
		"Instruction": c.GetCommandName(),
		"Host":        h.IP,
		"package":     packageName,
	}).Error(err.Error())
}
