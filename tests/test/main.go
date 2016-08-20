package main

import (
	"fmt"

	"log"

	"bufio"

	"golang.org/x/crypto/ssh"
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

	stdout, _ := session.StdoutPipe()

	/** IT WORKS*/
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println("Hello: " + scanner.Text()) // Println will add back the final '\n'
		}
	}()
	/* END IT WORKS */

	/** IT WORKS (more or less) */
	//reader := bufio.NewReader(stdout)
	//go func(){
	//	for {
	//		line, _ := reader.ReadString(byte('\n'))
	//		fmt.Printf("Hello: %s", line)
	//	}
	//}()
	/* END IT WORKS */

	session.Run("sudo yum install -y nano; sudo yum remove -y nano")

	var ss string
	fmt.Scanf(ss)
}
