package instructions

import (
	"io"
	"os"

	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/constants"
	"golang.org/x/crypto/ssh"

	log "github.com/Sirupsen/logrus"

	"strings"
)

type ENV struct {
	Name        string
	Description string
	Environment string
}

func (e *ENV) Execute(n connection.Node) {
	log.WithFields(log.Fields{
		"Instruction": "ENV",
		"Node":        n.IP,
	}).Info(constants.ARROW_LENGTH + e.Description)

	sshConfig := &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.Password),
		},
	}

	connection, err := ssh.Dial("tcp", n.IP+":22", sshConfig)
	if err != nil {
		log.Errorf("Failed to dial: %s\n", n.IP)
		log.Fatal(err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}

	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	env := strings.Split(e.Environment, "=")
	if len(env) == 2 {
		session.Setenv(env[0], env[1])
	} else {
		log.Fatal("More than one '=' found on ENV instruction")
	}
}
