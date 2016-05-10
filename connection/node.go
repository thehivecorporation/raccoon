package connection

import (
	"io"
	"os"

	"bytes"
	"errors"

	"fmt"

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
	Session      *ssh.Session `json:"session, omitempy"`
}

//GetSession will create a new connection to a Node and return it
func (n *Node) GetSession() (*ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.Password),
		},
	}

	connection, err := ssh.Dial("tcp", n.IP+":22", sshConfig)
	if err != nil {
		return nil, errors.New("Failed to dial: " + n.IP)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, errors.New("Failed to create session: " + err.Error())
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, errors.New("request for pseudo terminal failed: " + err.Error())
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, errors.New("Unable to setup stdout for session: " + err.Error())
	}

	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, errors.New("Unable to setup stderr for session: " + err.Error())
	}
	go io.Copy(os.Stderr, stderr)

	n.Session = session

	return session, nil
}

//CloseSession closes stored ssh session in Node. Remember to call it explicitly
//after all instructions has finished
func (n *Node) CloseSession() error {
	return n.Session.Close()
}

func (n *Node) GetOutput() (*ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: "vagrant",
		Auth: []ssh.AuthMethod{
			ssh.Password("vagrant"),
		},
	}

	connection, err := ssh.Dial("tcp", "192.168.31.10:22", sshConfig)
	if err != nil {
		return nil, errors.New("Failed to dial: 192.168.33.10")
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, errors.New("Failed to create session: " + err.Error())
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, errors.New("request for pseudo terminal failed: " + err.Error())
	}

	c := make(chan []byte, 10)

	go func(c *chan []byte) {
		for {
			m := <-*c
			println("Hello:" + string(m))
		}
	}(&c)

	r, w, _ := os.Pipe()
	r2, w2, _ := os.Pipe()

	session.Stdout = w
	session.Stderr = w2

	var byt bytes.Buffer
	var byt2 bytes.Buffer

	go func(c *chan []byte) {
		byt.ReadFrom(r)
		*c <- byt.Bytes()
	}(&c)

	go func(c *chan []byte) {
		byt2.ReadFrom(r2)
		*c <- byt2.Bytes()
	}(&c)

	//go io.Copy(os.Stdout, stdout)


	session.Run("ls /")

	var ss string
	fmt.Scanf(ss)
	return nil, nil
}
