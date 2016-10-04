package instructions

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
)

func logError(err error, c raccoon.CommandsExecutor, h *raccoon.Host) {
	log.WithFields(log.Fields{
		"command": c.GetCommand().Name,
		"host":    h.IP,
		"package": packageName,
		"color":    h.Color,
		"username": h.Username,
		"ssh_port": h.SSHPort,
	}).Error(err.Error())
}

func logCommand(fields map[string]interface{}, h raccoon.Host, c raccoon.CommandsExecutor) {
	commonFields := log.Fields{
		"command":  c.GetCommand().Name,
		"host":     h.IP,
		"package":  packageName,
		"username": h.Username,
		"ssh_port": h.SSHPort,
		"color":    h.Color,
	}

	if fields != nil {
		for k, field := range fields {
			commonFields[k] = field
		}
	}

	h.HostLogger.WithFields(commonFields).Info("> " + c.GetCommand().Description)
}
