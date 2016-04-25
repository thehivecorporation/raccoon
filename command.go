package main

import (
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func ExecuteCommandOnNode(c Command, n Node) {
	sshConfig := &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.Password),
		},
	}

	connection, err := ssh.Dial("tcp", n.IP+":22", sshConfig)
	if err != nil {
		log.Fatal("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatal("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatal("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	log.Printf("Command: %s\n", c.Name)
	err = session.Run(c.Command)
	if err != nil {
		log.Fatal(err)
	}
}
