package connection

import (
	"bufio"

	"github.com/Sirupsen/logrus"

	"errors"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"os"
)

//Node is a remote machine (virtual or physical) that we will execute our
//instructions on.
type Node struct {
	//IP of the remote host
	IP           string `json:"ip"`

	//Username to access remote host
	Username     string `json:"username,omitempty"`

	//Password to access remote host
	Password     string `json:"password,omitempty"`

	// TODO AuthFilePath corresponds to the path of the private key that could
	// give access to a remote machine.
	AuthFilePath string `json:"authFilePath,omitempty"`

	//Color that this node will output when printing in stdout
	Color color.Attribute

	sshClient *ssh.Client
}

//Specific logger for Node package
var log = logrus.New()

//Iterator used as "global variable" to get a new color from following array
var colorIter int = 0

//Contains all possible colors that could be used for logging
var colors [13]color.Attribute

func init() {
	//Sets our custom text formatter
	log.Formatter = new(NodeTextFormatter)

	// Output to stderr instead of stdout, could also be a file.
	log.Out = os.Stdout

	logrus.SetLevel(logrus.DebugLevel)

	//Fills the colors array
	colors = [13]color.Attribute{color.FgGreen, color.FgCyan, color.FgMagenta,
		color.FgYellow, color.FgHiCyan, color.FgHiBlue, color.FgHiGreen,
		color.FgHiMagenta, color.FgHiRed, color.FgHiYellow, color.FgBlue,
		color.FgRed, color.FgYellow}
}

//generateUniqueColor must be called every time a new node is created to assign
//a new color profile for the node
func (n *Node) generateUniqueColor() {
	if colorIter == 13 {
		colorIter = 0
	}

	n.Color = colors[colorIter]
	colorIter++
}

//InitializeNode must be called prior any execution on Nodes
func (n *Node) InitializeNode() error {
	n.generateUniqueColor()

	_, err := n.GetClient()
	if err != nil {
		return err
	}

	return nil
}

//InitializeNode will create a random color for out and a new connection to
// a Node returning it
func (n *Node) GetClient() (*ssh.Client, error) {

	log.WithFields(logrus.Fields{
		"host":     n.IP,
		"username": n.Username,
		"package":  "connection",
		"color":    n.Color,
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

//GetSession returns a new session once a client is created and the connection
//has been established successfully
func (n *Node) GetSession() (*ssh.Session, error) {
	if n.sshClient == nil {
		_, err := n.GetClient()

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
	if err != nil {
		return &ssh.Session{}, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return &ssh.Session{}, err
	}

	go n.sessionListenerRoutine(bufio.NewScanner(stdout))
	go n.sessionListenerRoutine(bufio.NewScanner(stderr))

	return session, nil
}

//sessionListenerRoutine is expected to use as a goroutine and receive as a
//parameter a scanner instance that is connected to an output
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
