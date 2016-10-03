package instructions

import (
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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

	if isFolder(a.SourcePath){
		files := pathsFromFilesInDir(path.Dir(a.SourcePath))
		for _, file := range files {
			add := ADD{
				SourcePath:file,
				DestPath:a.DestPath,
				Command:a.Command,
			}

			add.copyFile(h)
		}
	} else {
		a.copyFile(h)
	}
}

func (a *ADD) copyFile(h raccoon.Host) error {

	session, err := h.GetSession()
	if err != nil {
		logError(err, a, &h)
		return
	}
	defer session.Close()

	logCommand(log.Fields{
		"SourcePath": a.SourcePath,
		"DestPath":   a.DestPath,
	}, h, a)

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
			errorCh <- err
		}
		defer w.Close()

		if _, err := fmt.Fprintf(w, "C0644 %d %s\n", fi.Size(), fi.Name()); err != nil {
			errorCh <- err
			return
		}

		if _, err := io.Copy(w, f); err != nil {
			errorCh <- err
			return
		}

		if _, err := fmt.Fprint(w, "\x00"); err != nil { // transfer end with \x00
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

func pathsFromFilesInDir(d string)[]string{
	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Errorf("Error retrieving files from folder")
	}
	filePaths := make([]string, 0)
	for _, f := range files {
		filePaths = append(filePaths, f.Name())
	}
	return filePaths
}

func isFolder(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Errorf("Error retrieving info from file in ADD")
		return false
	}

	return fileInfo.IsDir()
}