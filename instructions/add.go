package instructions

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"path"
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
		if err := a.createFolder(h); err != nil {
			logError(err, a, &h)
		}

		files := pathsFromFilesInDir(a.SourcePath)
		for _, file := range files {
			add := ADD{
				SourcePath: path.Clean(a.SourcePath) + "/" + file,
				DestPath:path.Clean(a.DestPath) + "/" + path.Base(a.SourcePath),
				Command:a.Command,
			}

			err := add.copyFile(h)
			if err != nil {
				logError(err, a, &h)
			}
		}
	} else {
		a.copyFile(h)
	}
}

func (a *ADD) createFolder(h raccoon.Host) error {
	session, err := h.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	logCommand(nil, h, a)

	if err = session.Run("mkdir -p " + path.Clean(a.DestPath) + "/" + path.Base(a.SourcePath)); err != nil {
		return nil
	}

	return nil
}

func (a *ADD) copyFile(h raccoon.Host) error {

	session, err := h.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	logCommand(log.Fields{
		"SourcePath": a.SourcePath,
		"DestPath":   a.DestPath,
	}, h, a)

	f, err := os.Open(a.SourcePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = copyToSession(session, a.DestPath, f); err != nil {
		return err
	}

	return nil
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
		log.Errorf("Error retrieving info from file %s\n", path)
		return false
	}

	return fileInfo.IsDir()
}