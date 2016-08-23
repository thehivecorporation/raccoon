package instructions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

func logError(err error, c raccoon.CommandsExecutor, h *raccoon.Host) {
	log.WithFields(log.Fields{
		"Instruction": c.GetCommand().Name,
		"Host":        h.IP,
		"package":     packageName,
	}).Error(err.Error())
}

func logCommand(fields map[string]interface{}, ip string, c raccoon.CommandsExecutor) {
	commonFields := log.Fields{
		"Instruction": c.GetCommand().Name,
		"Node":        ip,
		"package":     packageName,
	}

	if fields != nil {
		for k, field := range fields {
			commonFields[k] = field
		}
	}

	log.WithFields(commonFields).Info(c.GetCommand().Description)
}
