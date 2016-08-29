package instructions

import (
	"testing"
	"time"

	"github.com/thehivecorporation/raccoon"
)

func TestExecuteAddInstruction(t *testing.T) {
	node := raccoon.Host{
		IP:       "127.0.0.1",
		Username: "root",
		Password: "root",
	}

	testContents := "Hello raccoon2\n"
	f, err := createTestFileWithContent(testContents)
	passOrError(err, t.Fatal)

	defer f.Close()

	add := ADD{
		Command: raccoon.Command{
			Name:        "ADD",
			Description: "Add description",
		},
		DestPath:    "/tmp",
		SourcePath:  f.Name(),
	}

	imageName := "rastasheep/ubuntu-sshd:14.04"
	containerID, err := launchContainer(imageName)
	passOrError(err, t.Fatal)

	defer cleanupContainer(containerID)

	time.Sleep(2 * time.Second)

	err = node.InitializeNode()
	passOrError(err, t.Fatal)

	add.Execute(node)

	testSession, err := createSessionPrototype()
	passOrError(err, t.Error)
	defer testSession.Close()

	testWriter := &IOWriterTester{
		T:              t,
		ExpectedString: testContents,
		Log:            true,
	}
	testSession.Stdout = testWriter
	testSession.Stderr = testWriter

	testSession.Run("cat " + f.Name())

	//Pass a non existing file
	add.SourcePath = "/tmp/iDoNotExist"
	add.Execute(node)
}

func TestCopyFileToHost(t *testing.T) {
	node := raccoon.Host{}
	node.IP = "127.0.0.1"
	node.Username = "root"
	node.Password = "root"

	add := ADD{
		Command: raccoon.Command{
			Name:        "ADD",
			Description: "Adding a file to destination target",
		},
		DestPath: "/tmp",
	}

	//Prepare: Launch a docker container with sshd running, user root and
	//password root. The file will be sent to this container
	imageName := "rastasheep/ubuntu-sshd:14.04"

	containerID, err := launchContainer(imageName)
	passOrError(err, t.Fatal)

	defer cleanupContainer(containerID)

	err = node.InitializeNode()
	passOrError(err, t.Fatal)

	//Prepare: Create a temp file to pass in with the contents "Hello Raccoon"
	//(without quotes)
	fileTestcontent := "Hello raccoon\n"
	f, err := createTestFileWithContent(fileTestcontent)
	passOrError(err, t.Fatal)

	add.SourcePath = f.Name()
	defer f.Close()

	//Prepare: Create a session with the container
	time.Sleep(2 * time.Second) //Give the container some time to launch
	session, err := node.GetSession()
	passOrError(err, t.Fatal)

	defer session.Close()

	//Prepare: Send the file with the information we have until now to /tmp
	err = copyToSession(session, add.DestPath, f)
	passOrError(err, t.Error)

	//Test: Check that the file exists in the target machine on path /tmp. The
	//file must have the same name and content in the target machine
	testSession, err := createSessionPrototype()
	passOrError(err, t.Fatal)

	defer testSession.Close()

	testWriter := &IOWriterTester{
		T:              t,
		ExpectedString: fileTestcontent,
		Log:            true,
	}
	testSession.Stdout = testWriter
	testSession.Stderr = testWriter

	testSession.Run("cat " + f.Name())
}
