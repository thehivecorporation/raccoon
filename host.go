package raccoon

import (
	"bufio"
	"fmt"

	"errors"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

//Host is a remote machine (virtual or physical) where we will execute our
//instructions on.
type Host struct {
	//IP of the remote host
	IP string `json:"ip"`

	//Username to access remote host
	Username string `json:"username,omitempty"`

	//Password to access remote host
	Password string `json:"password,omitempty"`

	// TODO AuthFilePath corresponds to the path of the private key that could
	// give access to a remote machine.
	AuthFilePath string `json:"authFilePath,omitempty"`

	//Color that this host will output when printing in stdout
	Color color.Attribute

	sshClient *ssh.Client
}

//Specific logger for Node package
var hostLogger = logrus.New()

//Iterator used as "global variable" to get a new color from following array
var colorIter int = 0

//Contains all possible colors that could be used for logging
var colors [13]color.Attribute

func init() {
	//Sets our custom text formatter
	hostLogger.Formatter = new(hostStdoutFormatter)

	// Output to stderr instead of stdout, could also be a file.
	hostLogger.Out = os.Stdout

	logrus.SetLevel(logrus.DebugLevel)

	//Fills the colors array
	colors = [13]color.Attribute{color.FgGreen, color.FgCyan, color.FgMagenta,
		color.FgYellow, color.FgHiCyan, color.FgHiBlue, color.FgHiGreen,
		color.FgHiMagenta, color.FgHiRed, color.FgHiYellow, color.FgBlue,
		color.FgRed, color.FgYellow}
}

//generateUniqueColor must be called every time a new host is created to assign
//a new color profile for the host
func (n *Host) generateUniqueColor() {
	if colorIter == 13 {
		colorIter = 0
	}

	n.Color = colors[colorIter]
	colorIter++
}

//InitializeNode must be called prior any execution on Hosts. It will try to
//get a ssh.Client object to open ssh session on host
func (n *Host) InitializeNode() error {
	n.generateUniqueColor()

	_, err := n.GetClient()
	if err != nil {
		return err
	}

	return nil
}

//GetClient will create a random color for logging and a new connection to
//a Host on port 22.
//
//TODO Make connection port configurable on JSON
func (n *Host) GetClient() (*ssh.Client, error) {

	hostLogger.WithFields(logrus.Fields{
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

//GetSession returns a new session once a client is created. In case the
//connection was lost, it tries to connect again with the Host
func (n *Host) GetSession() (*ssh.Session, error) {
	if n.sshClient == nil {
		_, err := n.GetClient()

		if err != nil {
			hostLogger.WithFields(logrus.Fields{
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

	stdout, err := session.StdoutPipe()
	if err != nil {
		return &ssh.Session{}, fmt.Errorf("Error getting stdout pipe from session: %s", err.Error())
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return &ssh.Session{}, fmt.Errorf("Error getting stdin pipe from session: %s", err.Error())
	}

	go n.sessionListenerRoutine(bufio.NewScanner(stdout))
	go n.sessionListenerRoutine(bufio.NewScanner(stderr))

	return session, nil
}

//sessionListenerRoutine is expected to use as a goroutine and receive as a
//parameter a scanner instance that is connected to an output
func (n *Host) sessionListenerRoutine(s *bufio.Scanner) {
	for s.Scan() {
		hostLogger.WithFields(logrus.Fields{
			"host":     n.IP,
			"username": n.Username,
			"package":  "connection",
			"color":    n.Color,
		}).Info(s.Text())
	}
}

//CloseNode closes stored ssh session in Host. Remember to call it explicitly
//after all instructions has finished
func (n *Host) CloseNode() error {
	if n.sshClient != nil {
		err := n.sshClient.Close()

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("ssh client was nil. Client was already closed")
}
