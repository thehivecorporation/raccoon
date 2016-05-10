package instructions

import (
	log "github.com/Sirupsen/logrus"

	"os"

	"io"

	"github.com/thehivecorporation/raccoon/connection"
	"golang.org/x/crypto/ssh"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/kubernetes/kubernetes/third_party/golang/go/doc/testdata"
	"errors"
)

//RUN is a instruction that in the recipe file correspond to the CMD instruction.
//It will execute the "Command" on every machine. Ideally, every command must
//be bash
type RUN struct {
	Name        string
	Description string
	Instruction string
}

//Execute is the implementation of the Instruction interface for a CMD instruction
func (c *RUN) Execute(n connection.Node) {
	log.WithFields(log.Fields{
		"Instruction": "RUN",
		"Node":        n.IP,
	}).Info(constants.ARROW_LENGTH + c.Description)

	//executeCommandOnNode(*c, n)

	err := executeShellCommandOnNode(*c, n)
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "RUN",
			"Node":        n.IP,
		}).Error(constants.ARROW_LENGTH + err.Error())
	}
}

func executeShellCommandOnNode(c RUN, n connection.Node) error {
	if n.Session != nil {
		err := n.Session.Run(c.Instruction)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Session was nil for " + n.IP)
}

func executeCommandOnNode(c RUN, n connection.Node) {
	log.WithFields(log.Fields{
		"Instruction": "RUN",
		"Node":        n.IP,
	}).Info(constants.ARROW_LENGTH + c.Description)

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

	session.Run(c.Instruction)
}
