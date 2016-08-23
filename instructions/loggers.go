package instructions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

func logError(err error, c raccoon.CommandsExecutor, h *raccoon.Host) {
	log.WithFields(log.Fields{
		"Instruction": c.GetCommandName(),
		"Host":        h.IP,
		"package":     packageName,
	}).Error(err.Error())
}

func logCommand(fields map[string]interface{}, ip, desc, command string) {
	commonFields := log.Fields{
		"Instruction": command,
		"Node":        ip,
		"package":     packageName,
	}

	if fields != nil {
		for k, field := range fields {
			commonFields[k] = field
		}
	}

	log.WithFields(commonFields).Info(desc)
}
