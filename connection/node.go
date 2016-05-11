package connection

import (
	"bufio"

	"github.com/Sirupsen/logrus"

	"errors"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

//Node is a remote machine (virtual or physical) that we will execute our
//instructions on. Fields are self descriptive. AuthFilePath corresponds to the
//path of the private key that could give access to a remote machine (TODO)
type Node struct {
	IP           string `json:"ip"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	AuthFilePath string `json:"authFilePath,omitempty"`

	sshClient *ssh.Client

	Color color.Attribute
}

var log = logrus.New()

var colorIter int = 0
var colors [13]color.Attribute

func init() {
	log.Formatter = new(NodeTextFormatter)

	logrus.SetLevel(logrus.DebugLevel)

	colors = [13]color.Attribute{color.FgHiCyan, color.FgBlue, color.FgYellow,
		color.FgCyan, color.FgGreen, color.FgHiBlue, color.FgHiGreen,
		color.FgHiMagenta, color.FgHiRed, color.FgHiYellow, color.FgMagenta,
		color.FgRed, color.FgYellow}
}

func (n *Node) GenerateUniqueColor() {
	if colorIter == 13 {
		colorIter = 0
	}

	n.Color = colors[colorIter]
	colorIter++
}

//InitializeNode will create a random color for out and a new connection to
// a Node returning it
func (n *Node) InitializeNode() (*ssh.Client, error) {
	log.WithFields(logrus.Fields{
		"host":     n.IP,
		"username": n.Username,
		"package":  "connection",
		"color":n.Color,
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

	n.sshClient = client

	return client, nil
}

func (n *Node) GetSession() (*ssh.Session, error) {
	if n.sshClient == nil {
		_, err := n.InitializeNode()

		if err != nil {
			log.WithFields(logrus.Fields{
				"host":    n.IP,
				"package": "dispatcher",
			}).Fatal("Error getting session: " + err.Error())

			return &ssh.Session{}, errors.New("Error getting session: " + err.Error())
		}
	}

	session, err := n.sshClient.NewSession()
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
	for s.Scan() {
		log.WithFields(logrus.Fields{
			"host":     n.IP,
			"username": n.Username,
			"package":  "connection",
			"color":    n.Color,
		}).Info(s.Text())
	}
}

//CloseSession closes stored ssh session in Node. Remember to call it explicitly
//after all instructions has finished
func (n *Node) CloseNode() error {
	if n.sshClient != nil {
		err := n.sshClient.Close()

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("ssh client was nil. Client was already closed")
}
