package main

import (
	"fmt"
	"os"

	"log"

	"golang.org/x/crypto/ssh"
	"gopkg.in/bufio.v1"
	"io"
)

func main() {
	sshConfig := &ssh.ClientConfig{
		User: "vagrant",
		Auth: []ssh.AuthMethod{
			ssh.Password("vagrant"),
		},
	}

	connection, err := ssh.Dial("tcp", "192.168.31.10:22", sshConfig)
	if err != nil {
		log.Fatal("Failed to dial: 192.168.33.10")
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatal("request for pseudo terminal failed: " + err.Error())
	}

	//stdout, _ := session.StdoutPipe()

	buf := bufio.NewBufferString("Hello: ")
	//go buf.WriteTo(os.Stdout)
	r, w := io.Pipe()
	session.Stdout = w
	go func(){
		for {
			buf.ReadFrom(r)
		}
	}()

	for {
		buf.WriteTo(os.Stdout)
	}
	//byt := make([]byte, 8)
	//go io.CopyBuffer(os.Stdout, stdout, byt)

	session.Run("ls /")

	var ss string
	fmt.Scanf(ss)

}
