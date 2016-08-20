package instructions

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"golang.org/x/crypto/ssh"
)

//ADD is a instruction to copy files with scp
type ADD struct {
	//Source path of the file
	SourcePath string
	//Destination path of the file
	DestPath string
	//Description of the instruction
	Description string
	//The name that identifies this struct ("ADD" in this case)
	Name string
}

//Execute is the implementation of the Instruction interface for a ADD instruction
func (c *ADD) Execute(n raccoon.Node) {
	session, err := n.GetSession()
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "ADD",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error(err.Error())
		session.Close()

		return
	}
	defer session.Close()

	log.WithFields(log.Fields{
		"Instruction": "ADD",
		"Node":        n.IP,
		"SourcePath":  c.SourcePath,
		"DestPath":    c.DestPath,
		"package":     "instructions",
	}).Info(c.Description)

	f, err := os.Open(c.SourcePath)
	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "ADD",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error(err.Error())
		return
	}
	defer f.Close()

	err = copyToSession(session, c.DestPath, f)

	if err != nil {
		log.WithFields(log.Fields{
			"Instruction": "ADD",
			"Node":        n.IP,
			"package":     "instructions",
		}).Error(err.Error())
	}
}

func copyToSession(session *ssh.Session, destinationFolder string, f *os.File) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	errorCh := make(chan error, 5)

	go func(errorCh chan error) {
		w, err := session.StdinPipe()
		if err != nil {
			println("error")
			errorCh <- err
		}
		defer w.Close()

		if _, err := fmt.Fprintf(w, "C0644 %d %s\n", fi.Size(), fi.Name()); err != nil {
			println("error")
			errorCh <- err
			return
		}

		if _, err := io.Copy(w, f); err != nil {
			println("error")
			errorCh <- err
			return
		}

		if _, err := fmt.Fprint(w, "\x00"); err != nil { // transfer end with \x00
			println("error")
			errorCh <- err
			return
		}

		errorCh <- nil

		close(errorCh)

	}(errorCh)

	if err = session.Run("/usr/bin/scp -tv " + destinationFolder); err != nil {
		return fmt.Errorf("Failed to run: %s" + err.Error())
	}

	for err := range errorCh {
		if err != nil {
			return err
		}
	}

	return nil
}
