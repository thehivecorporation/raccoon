package instructions

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"golang.org/x/crypto/ssh"
)

//ADD copies a single file to the destination host folder
type ADD struct {
	//Source path of the file
	SourcePath string
	//Destination path of the file in host
	DestPath string

	Command raccoon.Command
}

func (a *ADD) GetCommand() *raccoon.Command {
	return &a.Command
}

//Execute is the implementation of the Instruction interface for a ADD instruction
func (a *ADD) Execute(h raccoon.Host) {
	session, err := h.GetSession()
	if err != nil {
		logError(err, a, &h)
		return
	}
	defer session.Close()

	logCommand(log.Fields{
		"SourcePath": a.SourcePath,
		"DestPath":   a.DestPath,
	}, h.IP, a)

	f, err := os.Open(a.SourcePath)
	if err != nil {
		logError(err, a, &h)
		return
	}
	defer f.Close()

	if err = copyToSession(session, a.DestPath, f); err != nil {
		logError(err, a, &h)
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
