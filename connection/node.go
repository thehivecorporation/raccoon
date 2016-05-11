package connection

import (
	"bufio"

	"errors"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

//Node is a remote machine (virtual or physical) that we will execute our
//instructions on. Fields are self descriptive. AuthFilePath corresponds to the
//path of the private key that could give access to a remote machine (TODO)
type Node struct {
	IP           string       `json:"ip"`
	Username     string       `json:"username,omitempty"`
	Password     string       `json:"password,omitempty"`
	AuthFilePath string       `json:"authFilePath,omitempty"`

	client *ssh.Client
}



//GetSession will create a new connection to a Node and return it
func (n *Node) GetClient() (*ssh.Client, error) {
	log.WithFields(log.Fields{
		"host":     n.IP,
		"username": n.Username,
		"package":  "connection",
	}).Info("Opening SSH session")

	sshConfig := &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.Password),
		},
	}

	client, err := ssh.Dial("tcp", n.IP+":22", sshConfig)
	if err != nil {
		return nil, errors.New("Failed to dial: " + n.IP)
	}

	n.client = client

	return client, nil
}

func (n *Node) GetSession() (*ssh.Session, error) {
	session, err := n.client.NewSession()
	if err != nil {
		return &ssh.Session{}, errors.New("Failed to create session: " + err.Error())
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return &ssh.Session{}, errors.New("request for pseudo terminal failed: " + err.Error())
	}

	stdout, err := session.StdoutPipe()
	stderr, err := session.StderrPipe()

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go n.sessionListenerRoutine(stdoutScanner)
	go n.sessionListenerRoutine(stderrScanner)

	return session, nil
}

func (n *Node) sessionListenerRoutine(s *bufio.Scanner) {
	log.SetLevel(log.WarnLevel)

	for s.Scan() {
		// Println will add back the final '\n'
		log.WithFields(log.Fields{
			"host":     n.IP,
			"username": n.Username,
			"package":  "connection",
		}).Info("Message:" + s.Text())
	}
}

//CloseSession closes stored ssh session in Node. Remember to call it explicitly
//after all instructions has finished
func (n *Node) CloseClient() error {
	err := n.client.Close()

	if err != nil {
		return err
	}

	return nil
}
